// Package qml offers graphical QML application support for the Go language.
//
// Warning
//
// This package is in an alpha stage, and still in heavy development. APIs may
// change, and things may break.
//
// At this time contributors and developers that are interested in tracking the
// development closely are encouraged to use it. If you'd prefer a more stable
// release, please hold on a bit and subscribe to the mailing list for news. It's
// in a pretty good state, so it shall not take too long.
//
// See http://github.com/niemeyer/qml for details.
//
package qml

// #include <stdlib.h>
//
// #include "capi.h"
//
import "C"

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

// InitOptions holds options to initialize the qml package.
type InitOptions struct {
	// Reserved for coming options.
}

var initialized int32

// Init initializes the qml package with the provided parameters.
// If the options parameter is nil, default options suitable for a
// normal graphic application will be used.
//
// Init must be called only once, and before any other functionality
// from the qml package is used.
func Init(options *InitOptions) {
	if !atomic.CompareAndSwapInt32(&initialized, 0, 1) {
		panic("qml.Init called more than once")
	}

	guiLoopReady.Lock()
	go guiLoop()
	guiLoopReady.Lock()
}

// Engine provides an environment for instantiating QML components.
type Engine struct {
	Common
	values    map[interface{}]*valueFold
	destroyed bool

	imageProviders map[string]*func(providerId string, width, height int) image.Image
}

var engines = make(map[unsafe.Pointer]*Engine)

// NewEngine returns a new QML engine.
//
// The Destory method must be called to finalize the engine and
// release any resources used.
func NewEngine() *Engine {
	engine := &Engine{values: make(map[interface{}]*valueFold)}
	gui(func() {
		engine.addr = C.newEngine(nil)
		engine.engine = engine
		engine.imageProviders = make(map[string]*func(providerId string, width, height int) image.Image)
		engines[engine.addr] = engine
		stats.enginesAlive(+1)
	})
	return engine
}

func (e *Engine) assertValid() {
	if e.destroyed {
		panic("engine already destroyed")
	}
}

// Destroy finalizes the engine and releases any resources used.
// The engine must not be used after calling this method.
//
// It is safe to call Destroy more than once.
func (e *Engine) Destroy() {
	if !e.destroyed {
		gui(func() {
			if !e.destroyed {
				e.destroyed = true
				C.delObjectLater(e.addr)
				if len(e.values) == 0 {
					delete(engines, e.addr)
				} else {
					// The engine reference keeps those values alive.
					// The last value destroyed will clear it.
				}
				stats.enginesAlive(-1)
			}
		})
	}
}

// Load loads a new component with the provided location and with the
// content read from r. The location informs the resource name for
// logged messages, and its path is used to locate any other resources
// referenced by the QML content.
//
// Once a component is loaded, component instances may be created from
// the resulting object via its Create and CreateWindow methods.
func (e *Engine) Load(location string, r io.Reader) (Object, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if colon, slash := strings.Index(location, ":"), strings.Index(location, "/"); colon == -1 || slash <= colon {
		// TODO Better testing for this.
		if filepath.IsAbs(location) {
			location = "file:" + filepath.ToSlash(location)
		} else {
			dir, err := os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("cannot obtain absolute path: %v", err)
			}
			location = "file:" + filepath.ToSlash(filepath.Join(dir, location))
		}
	}

	cdata, cdatalen := unsafeBytesData(data)
	cloc, cloclen := unsafeStringData(location)
	comp := &Common{engine: e}
	gui(func() {
		// TODO The component's parent should probably be the engine.
		comp.addr = C.newComponent(e.addr, nilPtr)
		C.componentSetData(comp.addr, cdata, cdatalen, cloc, cloclen)
		message := C.componentErrorString(comp.addr)
		if message != nilCharPtr {
			err = errors.New(strings.TrimRight(C.GoString(message), "\n"))
			C.free(unsafe.Pointer(message))
		}
	})
	if err != nil {
		return nil, err
	}
	return comp, nil
}

// LoadFile loads a component from the provided QML file.
// Resources referenced by the QML content will be resolved relative to its path.
//
// Once a component is loaded, component instances may be created from
// the resulting object via its Create and CreateWindow methods.
func (e *Engine) LoadFile(path string) (Object, error) {
	// TODO Test this.
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return e.Load(path, f)
}

// LoadString loads a component from the provided QML string.
// The location informs the resource name for logged messages, and its
// path is used to locate any other resources referenced by the QML content.
//
// Once a component is loaded, component instances may be created from
// the resulting object via its Create and CreateWindow methods.
func (e *Engine) LoadString(location, qml string) (Object, error) {
	return e.Load(location, strings.NewReader(qml))
}

// Context returns the engine's root context.
func (e *Engine) Context() *Context {
	e.assertValid()
	var ctx Context
	ctx.engine = e
	gui(func() {
		ctx.addr = C.engineRootContext(e.addr)
	})
	return &ctx
}

// AddImageProvider registers f to be called when an image is requested by QML code
// with the specified provider identifier. It is a runtime error to register the same
// provider identifier multiple times.
//
// The imageId provided to f is the requested image source, with the "image:" scheme
// and provider identifier removed. For example, with an image image source of
// "image://myprovider/icons/home.ext", the respective imageId would be "icons/home.ext".
//
// If either the width or the height parameters provided to f are zero, no specific
// size for the image was requested. If non-zero, the returned image should have the
// the provided size, and will be resized if the returned image has a different size.
//
// See the documentation for more details on image providers:
//
//   http://qt-project.org/doc/qt-5.0/qtquick/qquickimageprovider.html
//
func (e *Engine) AddImageProvider(providerId string, f func(imageId string, width, height int) image.Image) {
	if _, ok := e.imageProviders[providerId]; ok {
		panic(fmt.Sprintf("engine already has an image provider with id %q", providerId))
	}
	e.imageProviders[providerId] = &f
	cproviderId, cproviderIdLen := unsafeStringData(providerId)
	gui(func() {
		qproviderId := C.newString(cproviderId, cproviderIdLen)
		defer C.delString(qproviderId)
		C.engineAddImageProvider(e.addr, qproviderId, unsafe.Pointer(&f))
	})
}

//export hookRequestImage
func hookRequestImage(imageFunc unsafe.Pointer, cid *C.char, cidLen, cwidth, cheight C.int) unsafe.Pointer {
	// TODO Test this.

	f := *(*func(imageId string, width, height int) image.Image)(imageFunc)

	id := unsafeString(cid, cidLen)
	width := int(cwidth)
	height := int(cheight)

	img := f(id, width, height)

	var cimage unsafe.Pointer

	rect := img.Bounds()
	width = rect.Max.X - rect.Min.X
	height = rect.Max.Y - rect.Min.Y
	cimage = C.newImage(C.int(width), C.int(height))

	var cbits []byte
	cbitsh := (*reflect.SliceHeader)((unsafe.Pointer)(&cbits))
	cbitsh.Data = (uintptr)((unsafe.Pointer)(C.imageBits(cimage)))
	cbitsh.Len = width * height * 4 // RGBA
	cbitsh.Cap = cbitsh.Len

	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			*(*uint32)(unsafe.Pointer(&cbits[i])) = (a>>8)<<24 | (r>>8)<<16 | (g>>8)<<8 | (b >> 8)
			i += 4
		}
	}
	return cimage
}

// Context represents a QML context that can hold variables visible
// to logic running within it.
type Context struct {
	Common
}

// SetVar makes the provided value available as a variable with the
// given name for QML code executed within the c context.
//
// If value is a struct, its exported fields are also made accessible to
// QML code as attributes of the named object. The attribute name in the
// object has the same name of the Go field name, except for the first
// letter which is lowercased. This is conventional and enforced by
// the QML implementation.
//
// The engine will hold a reference to the provided value, so it will
// not be garbage collected until the engine is destroyed, even if the
// value is unused or changed.
func (ctx *Context) SetVar(name string, value interface{}) {
	cname, cnamelen := unsafeStringData(name)
	gui(func() {
		var dvalue C.DataValue
		packDataValue(value, &dvalue, ctx.engine, cppOwner)

		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)

		C.contextSetProperty(ctx.addr, qname, &dvalue)
	})
}

// SetVars makes the exported fields of the provided value available as
// variables for QML code executed within the c context. The variable names
// will have the same name of the Go field names, except for the first
// letter which is lowercased. This is conventional and enforced by
// the QML implementation.
//
// The engine will hold a reference to the provided value, so it will
// not be garbage collected until the engine is destroyed, even if the
// value is unused or changed.
func (ctx *Context) SetVars(value interface{}) {
	gui(func() {
		C.contextSetObject(ctx.addr, wrapGoValue(ctx.engine, value, cppOwner))
	})
}

// Var returns the context variable with the given name.
func (ctx *Context) Var(name string) interface{} {
	cname, cnamelen := unsafeStringData(name)

	var dvalue C.DataValue
	gui(func() {
		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)

		C.contextGetProperty(ctx.addr, qname, &dvalue)
	})
	return unpackDataValue(&dvalue, ctx.engine)
}

// TODO Context.Spawn() => Context

// TODO engine.ObjectOf(&value) => *Common for the Go value

// Object is the common interface implemented by all QML types.
//
// See the documentation of Common for details about this interface.
type Object interface {
	Common() *Common
	Set(property string, value interface{}) error
	Property(name string) interface{}
	Int(property string) int
	Int64(property string) int64
	Float64(property string) float64
	Bool(property string) bool
	String(property string) string
	Color(property string) color.RGBA
	Object(property string) Object
	ObjectByName(objectName string) Object
	Call(method string, params ...interface{}) interface{}
	Create(ctx *Context) Object
	CreateWindow(ctx *Context) *Window
	Destroy()
	On(signal string, function interface{})
}

// Common implements the common behavior of all QML objects.
// It implements the Object interface.
type Common struct {
	addr   unsafe.Pointer
	engine *Engine
}

var _ Object = (*Common)(nil)

// Common returns obj itself.
//
// This provides access to the underlying *Common for types that
// embed it, when these are used via the Object interface.
func (obj *Common) Common() *Common {
	return obj
}

// Set changes the named object property to the given value.
func (obj *Common) Set(property string, value interface{}) error {
	cproperty := C.CString(property)
	defer C.free(unsafe.Pointer(cproperty))
	gui(func() {
		var dvalue C.DataValue
		packDataValue(value, &dvalue, obj.engine, cppOwner)
		C.objectSetProperty(obj.addr, cproperty, &dvalue)
	})
	// TODO Return an error if the value cannot be set.
	return nil
}

// Property returns the current value for a property of the object.
// If the property type is known, type-specific methods such as Int
// and String are more convenient to use.
// Property panics if the property does not exist.
func (obj *Common) Property(name string) interface{} {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var dvalue C.DataValue
	var found C.int
	gui(func() {
		found = C.objectGetProperty(obj.addr, cname, &dvalue)
	})
	if found == 0 {
		panic(fmt.Sprintf("object does not have a %q property", name))
	}
	return unpackDataValue(&dvalue, obj.engine)
}

// Int returns the int value of the named property.
// Int panics if the property cannot be represented as an int.
func (obj *Common) Int(property string) int {
	switch value := obj.Property(property).(type) {
	case int:
		return value
	case int64:
		if int64(int(value)) != value {
			panic(fmt.Sprintf("value of property %q is too large for int: %#v", property, value))
		}
		return int(value)
	case float32:
		return int(value)
	case float64:
		return int(value)
	default:
		panic(fmt.Sprintf("value of property %q cannot be represented as an int: %#v", property, value))
	}
}

// Int64 returns the int64 value of the named property.
// Int64 panics if the property cannot be represented as an int64.
func (obj *Common) Int64(property string) int64 {
	switch value := obj.Property(property).(type) {
	case int:
		return int64(value)
	case int64:
		return value
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	default:
		panic(fmt.Sprintf("value of property %q cannot be represented as an int64: %#v", property, value))
	}
}

// Float64 returns the float64 value of the named property.
// Float64 panics if the property cannot be represented as float64.
func (obj *Common) Float64(property string) float64 {
	switch value := obj.Property(property).(type) {
	case int:
		return float64(value)
	case int64:
		return float64(value)
	case float32:
		return float64(value)
	case float64:
		return value
	default:
		panic(fmt.Sprintf("value of property %q cannot be represented as a float64: %#v", property, value))
	}
}

// Bool returns the bool value of the named property.
// Bool panics if the property is not a bool.
func (obj *Common) Bool(property string) bool {
	value := obj.Property(property)
	if b, ok := value.(bool); ok {
		return b
	}
	panic(fmt.Sprintf("value of property %q is not a bool: %#v", property, value))
}

// String returns the string value of the named property.
// String panics if the property is not a string.
func (obj *Common) String(property string) string {
	value := obj.Property(property)
	if s, ok := value.(string); ok {
		return s
	}
	panic(fmt.Sprintf("value of property %q is not a string: %#v", property, value))
}

// Color returns the RGBA value of the named property.
// Color panics if the property is not a color.
func (obj *Common) Color(property string) color.RGBA {
	value := obj.Property(property)
	c, ok := value.(color.RGBA)
	if !ok {
		panic(fmt.Sprintf("value of property %q is not a color: %#v", property, value))
	}
	return c
}

// Object returns the object value of the named property.
// Object panics if the property is not a QML object.
func (obj *Common) Object(property string) Object {
	value := obj.Property(property)
	object, ok := value.(Object)
	if !ok {
		panic(fmt.Sprintf("value of property %q is not a QML object: %#v", property, value))
	}
	return object
}

// ObjectByName returns the Object value of the descendant object that
// was defined with the objectName property set to the provided value.
// ObjectByName panics if the object is not found.
func (obj *Common) ObjectByName(objectName string) Object {
	cname, cnamelen := unsafeStringData(objectName)
	var dvalue C.DataValue
	gui(func() {
		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)
		C.objectFindChild(obj.addr, qname, &dvalue)
	})
	object, ok := unpackDataValue(&dvalue, obj.engine).(Object)
	if !ok {
		panic(fmt.Sprintf("cannot find descendant with objectName == %q", objectName))
	}
	return object
}

// Call calls the given object method with the provided parameters.
// Call panics if the method does not exist.
func (obj *Common) Call(method string, params ...interface{}) interface{} {
	if len(params) > len(dataValueArray) {
		panic("too many parameters")
	}
	cmethod := C.CString(method)
	defer C.free(unsafe.Pointer(cmethod))
	var result C.DataValue
	gui(func() {
		for i, param := range params {
			packDataValue(param, &dataValueArray[i], obj.engine, jsOwner)
		}
		// TODO Panic if the underlying invokation returns false.
		// TODO Is there any other actual error other than existence that can be observed?
		//      If so, this method needs an error result too.
		C.objectInvoke(obj.addr, cmethod, &result, &dataValueArray[0], C.int(len(params)))
	})
	return unpackDataValue(&result, obj.engine)
}

// Create creates a new instance of the component held by obj.
// The component instance runs under the ctx context. If ctx is nil,
// it runs under the same context as obj.
//
// The Create method panics if called on an object that does not
// represent a QML component.
func (obj *Common) Create(ctx *Context) Object {
	if C.objectIsComponent(obj.addr) == 0 {
		panic("object is not a component")
	}
	var root Common
	root.engine = obj.engine
	gui(func() {
		ctxaddr := nilPtr
		if ctx != nil {
			ctxaddr = ctx.addr
		}
		root.addr = C.componentCreate(obj.addr, ctxaddr)
	})
	return &root
}

// CreateWindow creates a new instance of the component held by obj,
// and creates a new window holding the instance as its root object.
// The component instance runs under the ctx context. If ctx is nil,
// it runs under the same context as obj.
//
// The CreateWindow method panics if called on an object that
// does not represent a QML component.
func (obj *Common) CreateWindow(ctx *Context) *Window {
	if C.objectIsComponent(obj.addr) == 0 {
		panic("object is not a component")
	}
	var win Window
	win.engine = obj.engine
	gui(func() {
		ctxaddr := nilPtr
		if ctx != nil {
			ctxaddr = ctx.addr
		}
		win.addr = C.componentCreateWindow(obj.addr, ctxaddr)
	})
	return &win
}

// Destroy finalizes the value and releases any resources used.
// The value must not be used after calling this method.
func (obj *Common) Destroy() {
	// TODO We might hook into the destroyed signal, and prevent this object
	//      from being used in post-destruction crash-prone ways.
	gui(func() {
		if obj.addr != nilPtr {
			C.delObjectLater(obj.addr)
			obj.addr = nilPtr
		}
	})
}

var connectedFunction = make(map[*interface{}]bool)

// On connects the named signal from obj with the provided function, so that
// when obj next emits that signal, the function is called with the parameters
// the signal carries.
//
// The provided function must accept a number of parameters that is equal to
// or less than the number of parameters provided by the signal, and the
// resepctive parameter types must match exactly or be conversible according
// to normal Go rules.
//
// For example:
//
//     obj.On("clicked", func() { fmt.Println("obj got a click") })
//
// Note that Go uses the real signal name, rather than the one used when
// defining QML signal handlers ("clicked" rather than "onClicked").
//
// For more details regarding signals and QML see:
//
//     http://qt-project.org/doc/qt-5.0/qtqml/qml-qtquick2-connections.html
//
func (obj *Common) On(signal string, function interface{}) {
	funcv := reflect.ValueOf(function)
	funct := funcv.Type()
	if funcv.Kind() != reflect.Func {
		panic("function provided to On is not a function or method")
	}
	if funct.NumIn() > C.MaxParams {
		panic("function takes too many arguments")
	}
	csignal, csignallen := unsafeStringData(signal)
	var cerr *C.error
	gui(func() {
		cerr = C.objectConnect(obj.addr, csignal, csignallen, obj.engine.addr, unsafe.Pointer(&function), C.int(funcv.Type().NumIn()))
		if cerr == nil {
			connectedFunction[&function] = true
			stats.connectionsAlive(+1)
		}
	})
	if cerr != nil {
		panic(cerror(cerr).Error())
	}
}

//export hookSignalDisconnect
func hookSignalDisconnect(funcp unsafe.Pointer) {
	before := len(connectedFunction)
	delete(connectedFunction, (*interface{})(funcp))
	if before == len(connectedFunction) {
		panic("disconnecting unknown signal function")
	}
	stats.connectionsAlive(-1)
}

//export hookSignalCall
func hookSignalCall(enginep unsafe.Pointer, funcp unsafe.Pointer, args *C.DataValue) {
	engine := engines[enginep]
	if engine == nil {
		panic("signal called after engine was destroyed")
	}
	funcv := reflect.ValueOf(*(*interface{})(funcp))
	funct := funcv.Type()
	numIn := funct.NumIn()
	var params [C.MaxParams]reflect.Value
	for i := 0; i < numIn; i++ {
		arg := (*C.DataValue)(unsafe.Pointer(uintptr(unsafe.Pointer(args)) + uintptr(i)*dataValueSize))
		param := reflect.ValueOf(unpackDataValue(arg, engine))
		if paramt := funct.In(i); param.Type() != paramt {
			// TODO Provide a better error message when this fails.
			param = param.Convert(paramt)
		}
		params[i] = param
	}
	funcv.Call(params[:numIn])
}

func cerror(cerr *C.error) error {
	err := errors.New(C.GoString((*C.char)(unsafe.Pointer(cerr))))
	C.free(unsafe.Pointer(cerr))
	return err
}

// TODO Signal emitting support for go values.

// Window represents a QML window where components are rendered.
type Window struct {
	Common
}

// Show exposes the window.
func (win *Window) Show() {
	gui(func() {
		C.windowShow(win.addr)
	})
}

// Hide hides the window.
func (win *Window) Hide() {
	gui(func() {
		C.windowHide(win.addr)
	})
}

// Root returns the root object being rendered.
//
// If the window was defined in QML code, the root object is the window itself.
func (win *Window) Root() Object {
	var obj Common
	obj.engine = win.engine
	gui(func() {
		obj.addr = C.windowRootObject(win.addr)
	})
	return &obj
}

// Wait blocks the current goroutine until the window is closed.
func (win *Window) Wait() {
	// XXX Test this.
	var m sync.Mutex
	m.Lock()
	gui(func() {
		// TODO Must be able to wait for the same Window from multiple goroutines.
		// TODO If the window is not visible, must return immediately.
		waitingWindows[win.addr] = &m
		C.windowConnectHidden(win.addr)
	})
	m.Lock()
}

// Snapshot returns an image with the visible contents of the window.
// The main GUI thread is paused while the data is being acquired.
func (win *Window) Snapshot() image.Image {
	// TODO Test this.
	var cimage unsafe.Pointer
	gui(func() {
		cimage = C.windowGrabWindow(win.addr)
	})
	defer C.delImage(cimage)

	// This should be safe to be done out of the main GUI thread.
	var cwidth, cheight C.int
	C.imageSize(cimage, &cwidth, &cheight)

	var cbits []byte
	cbitsh := (*reflect.SliceHeader)((unsafe.Pointer)(&cbits))
	cbitsh.Data = (uintptr)((unsafe.Pointer)(C.imageConstBits(cimage)))
	cbitsh.Len = int(cwidth * cheight * 8) // ARGB
	cbitsh.Cap = cbitsh.Len

	image := image.NewRGBA(image.Rect(0, 0, int(cwidth), int(cheight)))
	l := int(cwidth * cheight * 4)
	for i := 0; i < l; i += 4 {
		var c uint32 = *(*uint32)(unsafe.Pointer(&cbits[i]))
		image.Pix[i+0] = byte(c >> 16)
		image.Pix[i+1] = byte(c >> 8)
		image.Pix[i+2] = byte(c)
		image.Pix[i+3] = byte(c >> 24)
	}
	return image
}

var waitingWindows = make(map[unsafe.Pointer]*sync.Mutex)

//export hookWindowHidden
func hookWindowHidden(addr unsafe.Pointer) {
	m, ok := waitingWindows[addr]
	if !ok {
		panic("window is not waiting")
	}
	delete(waitingWindows, addr)
	m.Unlock()
}

type TypeSpec struct {
	Location     string
	Major, Minor int
	// TODO Consider refactoring this type into ModuleSpec for the above + []TypeSpec for the below
	Name string
	New  func() interface{}
}

var types []*TypeSpec

func RegisterType(spec *TypeSpec) error {
	return registerType(spec, false)
}

func RegisterSingleton(spec *TypeSpec) error {
	return registerType(spec, true)
}

func registerType(spec *TypeSpec, singleton bool) error {
	// Copy and hold a reference to the spec data.
	localSpec := *spec

	// TODO Validate localSpec fields.

	var err error
	gui(func() {
		sample := spec.New()
		if sample == nil {
			err = fmt.Errorf("TypeSpec.New for type %q returned nil", spec.Name)
			return
		}

		cloc := C.CString(localSpec.Location)
		cname := C.CString(localSpec.Name)
		cres := C.int(0)
		if singleton {
			cres = C.registerSingleton(cloc, C.int(localSpec.Major), C.int(localSpec.Minor), cname, typeInfo(sample), unsafe.Pointer(&localSpec))
		} else {
			cres = C.registerType(cloc, C.int(localSpec.Major), C.int(localSpec.Minor), cname, typeInfo(sample), unsafe.Pointer(&localSpec))
		}
		// It doesn't look like it keeps references to these, but it's undocumented and unclear.
		C.free(unsafe.Pointer(cloc))
		C.free(unsafe.Pointer(cname))
		if cres == -1 {
			err = fmt.Errorf("QML engine failed to register type; invalid type location or name?")
		} else {
			types = append(types, &localSpec)
		}
	})

	return err
}

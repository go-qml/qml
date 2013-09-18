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
	"io"
	"io/ioutil"
	"os"
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
	addr      unsafe.Pointer
	values    map[interface{}]*valueFold
	destroyed bool
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
func (e *Engine) Load(location string, r io.Reader) (*Component, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if colon, slash := strings.Index(location, ":"), strings.Index(location, "/"); slash <= colon {
		location = "file:" + location
	}

	cdata, cdatalen := unsafeBytesData(data)
	cloc, cloclen := unsafeStringData(location)
	component := &Component{Object{engine: e}}
	gui(func() {
		// TODO The component's parent should probably be the engine.
		component.obj.addr = C.newComponent(e.addr, nilPtr)
		C.componentSetData(component.obj.addr, cdata, cdatalen, cloc, cloclen)
		message := C.componentErrorString(component.obj.addr)
		if message != nilCharPtr {
			err = errors.New(strings.TrimRight(C.GoString(message), "\n"))
			C.free(unsafe.Pointer(message))
		}
	})
	if err != nil {
		return nil, err
	}
	return component, nil
}

// Load loads a component from the provided QML file.
// Resources referenced by the QML content will be resolved relative to its path.
func (e *Engine) LoadFile(path string) (*Component, error) {
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
func (e *Engine) LoadString(location, qml string) (*Component, error) {
	return e.Load(location, strings.NewReader(qml))
}

// Context returns the engine's root context.
func (e *Engine) Context() *Context {
	e.assertValid()
	var ctx Context
	ctx.obj.engine = e
	gui(func() {
		ctx.obj.addr = C.engineRootContext(e.addr)
	})
	return &ctx
}

// Context represents a QML context that can hold variables visible
// to logic running within it.
type Context struct {
	obj Object
}

// TODO Consider whether to expose the methods of Object directly
//      on Context and Window by embedding it, or whether to have an
//      AsObject method that returns it.

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
		packDataValue(value, &dvalue, ctx.obj.engine, cppOwner)

		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)

		C.contextSetProperty(ctx.obj.addr, qname, &dvalue)
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
		C.contextSetObject(ctx.obj.addr, wrapGoValue(ctx.obj.engine, value, cppOwner))
	})
}

// Var returns the context variable with the given name.
func (ctx *Context) Var(name string) interface{} {
	cname, cnamelen := unsafeStringData(name)

	var dvalue C.DataValue
	gui(func() {
		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)

		C.contextGetProperty(ctx.obj.addr, qname, &dvalue)
	})
	return unpackDataValue(&dvalue, ctx.obj.engine)
}

// TODO Context.Spawn() => Context

type Component struct {
	obj Object
}

// TODO Drop Component. We should be able to use a plain qml.Object for these features,
//      and it makes sense to do so given that components may be defined in QML as well.

// Create creates an instance of the component.
func (c *Component) Create(ctx *Context) *Object {
	var obj Object
	obj.engine = c.obj.engine
	gui(func() {
		ctxaddr := nilPtr
		if ctx != nil {
			ctxaddr = ctx.obj.addr
		}
		obj.addr = C.componentCreate(c.obj.addr, ctxaddr)
	})
	return &obj
}

// CreateWindow creates a new instance of the c component running under
// the provided context, and creates a new window for the component
// instance to render its content into.
//
// If the provided context is nil, the engine's root context is used.
//
// If the returned window is not destroyed explicitly, it will be
// destroyed when the engine behind the used context is.
func (c *Component) CreateWindow(ctx *Context) *Window {
	var win Window
	win.obj.engine = c.obj.engine
	gui(func() {
		ctxaddr := nilPtr
		if ctx != nil {
			ctxaddr = ctx.obj.addr
		}
		win.obj.addr = C.componentCreateView(c.obj.addr, ctxaddr)
	})
	return &win
}

// TODO engine.ObjectOf(&value) => *Object for the Go value

type Object struct {
	addr   unsafe.Pointer
	engine *Engine
}

// Set changes the value of property.
func (obj *Object) Set(property string, value interface{}) error {
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
func (obj *Object) Property(name string) interface{} {
	// TODO Return an ok bool indicating whether the property was found.
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var dvalue C.DataValue
	gui(func() {
		C.objectGetProperty(obj.addr, cname, &dvalue)
	})
	return unpackDataValue(&dvalue, obj.engine)
}

// Int returns the int value of the given property.
// The call panics if the property value cannot be represented as an int.
func (obj *Object) Int(property string) int {
	switch value := obj.Property(property).(type) {
	case int:
		return value
	case int32:
		return int(value)
	case int64:
		if int64(int(value)) != value {
			panic(fmt.Sprintf("value of property %q is too large for int: %v", property, value))
		}
		return int(value)
	case float64:
		// May truncate, but seems a bit too much computing to validate these all the time.
		return int(value)
	default:
		panic(fmt.Sprintf("value of property %q cannot be represented as an int: %v", property, value))
	}
}

// String returns the string value of the given property.
// The call panics if the property value is not a string.
func (obj *Object) String(property string) string {
	s, ok := obj.Property(property).(string)
	if !ok {
		panic(fmt.Sprintf("value of property %q is not a string: %v", property, s))
	}
	return s
}

// TODO More type-specific methods: int64, float64, etc

// Object returns the *qml.Object value of the given property.
// The call panics if the property value is not a *qml.Object.
func (obj *Object) Object(property string) *Object {
	object, ok := obj.Property(property).(*Object)
	if !ok {
		panic(fmt.Sprintf("value of property %q is not a *qml.Object: %v", property, object))
	}
	return object
}

// ObjectByName returns the *qml.Object value of the descendant object that
// was defined with the objectName property set to the provided value.
// The call panics if the object is not found.
func (obj *Object) ObjectByName(objectName string) *Object {
	cname, cnamelen := unsafeStringData(objectName)
	var dvalue C.DataValue
	gui(func() {
		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)
		C.objectFindChild(obj.addr, qname, &dvalue)
	})
	object, ok := unpackDataValue(&dvalue, obj.engine).(*Object)
	if !ok {
		panic(fmt.Sprintf("cannot find descendant with objectName == %q", objectName))
	}
	return object
}

// TODO Consider using a Result wrapper type to be used by the Object.Call,
//      Object.Property, and Context.Var methods. It would offer methods such as
//      Int, and String, to facilitate converting (rather than just type-asserting)
//      results to the desired types, in a way equivalent to what Object currently
//      does for properties.

// Call calls the given object method with the provided parameters.
// It panics if the method does not exist.
func (obj *Object) Call(method string, params ...interface{}) interface{} {
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

func (obj *Object) Create(ctx *Context) *Object {
	// TODO Implement C.objectIsComponent and panic if it returns false.
	var root Object
	root.engine = obj.engine
	gui(func() {
		ctxaddr := nilPtr
		if ctx != nil {
			ctxaddr = ctx.obj.addr
		}
		root.addr = C.componentCreate(obj.addr, ctxaddr)
	})
	return &root
}

// Destroy finalizes the value and releases any resources used.
// The value must not be used after calling this method.
func (obj *Object) Destroy() {
	// TODO We might hook into the destroyed signal, and prevent this object
	//      from being used in post-destruction crash-prone ways.
	gui(func() {
		if obj.addr != nilPtr {
			C.delObjectLater(obj.addr)
			obj.addr = nilPtr
		}
	})
}

// TODO Object.Connect(name, func(...) {})

// TODO Signal emitting support for go values.

// Window represents a QML window where components are rendered.
type Window struct {
	obj Object
}

// Show exposes the window.
func (win *Window) Show() {
	gui(func() {
		C.viewShow(win.obj.addr)
	})
}

// Hide hides the window.
func (win *Window) Hide() {
	gui(func() {
		C.viewHide(win.obj.addr)
	})
}

// Root returns the root object being rendered in the window.
func (win *Window) Root() *Object {
	var obj Object
	obj.engine = win.obj.engine
	gui(func() {
		obj.addr = C.viewRootObject(win.obj.addr)
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
		waitingWindows[win.obj.addr] = &m
		C.viewConnectHidden(win.obj.addr)
	})
	m.Lock()
}

// Destroy destroys the window.
// The window should not be used after this method is called.
func (win *Window) Destroy() {
	win.obj.Destroy()
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
		if singleton {
			C.registerSingleton(cloc, C.int(localSpec.Major), C.int(localSpec.Minor), cname, typeInfo(sample), unsafe.Pointer(&localSpec))
		} else {
			C.registerType(cloc, C.int(localSpec.Major), C.int(localSpec.Minor), cname, typeInfo(sample), unsafe.Pointer(&localSpec))
		}
		// TODO Check if qmlRegisterType keeps a reference to those.
		//C.free(unsafe.Pointer(cloc))
		//C.free(unsafe.Pointer(cname))
		types = append(types, &localSpec)
	})

	// TODO Are there really no errors possible from qmlRegisterType?
	return err
}

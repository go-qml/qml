package qml

// #include <stdlib.h>
//
// #include "capi.h"
//
// void hack(void *engine, void *component);
import "C"

import (
	"errors"
	"fmt"
	"io/ioutil"
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

func Hack(engine *Engine, component *Component) {
	gui(func() {
		C.hack(engine.addr, component.addr)
	})
}

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

	fmt.Println("main() thread:", C.currentThread())
	go guiLoop()

	// Wait for app to be created and event loop to be running.
	gui(func() {})
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
		fmt.Println("NewEngine thread:", C.currentThread())
		engine.addr = C.newEngine(nil)
		engines[engine.addr] = engine
		stats.enginesAlive(+1)
	})
	fmt.Println("main() thread:", C.currentThread())
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
				C.delObject(e.addr)
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

type Content interface {
	Location() string
	Data() ([]byte, error)
}

func String(location, qml string) Content {
	return &content{location, []byte(qml), nil}
}

func File(path string) Content {
	// TODO: Test this.
	data, err := ioutil.ReadFile(path)
	return &content{path, data, err}
}

type content struct {
	location string
	data     []byte
	err      error
}

func (c *content) Location() string {
	return c.location
}

func (c *content) Data() ([]byte, error) {
	return c.data, c.err
}

// Load loads a new component with the provided QML content.
//
// For example:
//
//     component, err := engine.Load(qml.File("file.qml"))
//
// See qml.File and qml.String.
func (e *Engine) Load(c Content) (*Component, error) {
	data, err := c.Data()
	if err != nil {
		return nil, err
	}
	return e.newComponent(c.Location(), data)
}

// Context returns the engine's root context.
func (e *Engine) Context() *Context {
	e.assertValid()
	var context Context
	context.engine = e
	gui(func() {
		context.addr = C.engineRootContext(e.addr)
	})
	return &context
}

type Context struct {
	commonObject
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
func (c *Context) SetVar(name string, value interface{}) {
	cname, cnamelen := unsafeStringData(name)
	gui(func() {
		var dvalue C.DataValue
		packDataValue(value, &dvalue, c.engine, cppOwner)

		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)

		C.contextSetProperty(c.addr, qname, &dvalue)
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
func (c *Context) SetVars(value interface{}) {
	gui(func() {
		C.contextSetObject(c.addr, wrapGoValue(c.engine, value, cppOwner))
	})
}

// Var returns the context variable with the given name.
func (c *Context) Var(name string) interface{} {
	cname, cnamelen := unsafeStringData(name)

	var dvalue C.DataValue
	gui(func() {
		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)

		C.contextGetProperty(c.addr, qname, &dvalue)
	})
	return unpackDataValue(&dvalue, c.engine)
}

// TODO Context.Spawn() => Context

type Component struct {
	addr   unsafe.Pointer
	engine *Engine
}

func (e *Engine) newComponent(location string, data []byte) (*Component, error) {
	// TODO What's a nice way to delete the component and created component objects?
	cdata, cdatalen := unsafeBytesData(data)
	cloc, cloclen := unsafeStringData(location)
	component := &Component{engine: e}
	var err error
	gui(func() {
		component.addr = C.newComponent(e.addr, nilPtr)
		_, _, _, _ = cdata, cdatalen, cloc, cloclen
		_ = errors.New
		_ = strings.Split
		//C.componentSetData(component.addr, cdata, cdatalen, cloc, cloclen)
		//message := C.componentErrorString(component.addr)
		//if message != nilCharPtr {
		//	err = errors.New(strings.TrimRight(C.GoString(message), "\n"))
		//	C.free(unsafe.Pointer(message))
		//}
	})
	if err != nil {
		return nil, err
	}
	return component, nil
}

func (c *Component) Create(context *Context) *Value {
	var value Value
	value.engine = c.engine
	gui(func() {
		value.addr = C.componentCreate(c.addr, context.addr)
	})
	return &value
}

// CreateWindow creates a new instance of the c component running under
// the provided context, and creates a new window for the component
// instance to render its content into.
//
// If the provided context is nil, the engine's root context is used.
//
// If the returned window is not destroyed explicitly, it will be
// destroyed when the engine behind the used context is.
func (c *Component) CreateWindow(context *Context) *Window {
	if context == nil {
		// TODO Test this.
		context = c.engine.Context()
	}
	var window Window
	window.engine = c.engine
	gui(func() {
		window.addr = C.componentCreateView(c.addr, context.addr)
	})
	return &window
}

type commonObject struct {
	addr   unsafe.Pointer
	engine *Engine
}

// TODO engine.ValueOf(&value) => *Value for the Go value

type Value struct {
	commonObject
}

func (o *commonObject) Field(name string) interface{} {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var dvalue C.DataValue
	gui(func() {
		C.objectGetProperty(o.addr, cname, &dvalue)
	})
	return unpackDataValue(&dvalue, o.engine)
}

func (o *commonObject) SetField(name string, value interface{}) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	gui(func() {
		var dvalue C.DataValue
		packDataValue(value, &dvalue, o.engine, cppOwner)
		// TODO Handle the return value.
		C.objectSetProperty(o.addr, cname, &dvalue)
	})
}

func (o *commonObject) MustFind(name string) *Value {
	cname, cnamelen := unsafeStringData(name)
	var dvalue C.DataValue
	gui(func() {
		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)
		C.objectFindChild(o.addr, qname, &dvalue)
	})
	value, ok := unpackDataValue(&dvalue, o.engine).(*Value)
	if !ok {
		panic(fmt.Sprintf("cannot find child %q", name))
	}
	return value
}

func (o *commonObject) Call(method string, params ...interface{}) interface{} {
	// TODO Return errors.
	if len(params) > len(dataValueArray) {
		panic("too many parameters")
	}
	cmethod := C.CString(method)
	defer C.free(unsafe.Pointer(cmethod))
	var result C.DataValue
	gui(func() {
		for i, param := range params {
			packDataValue(param, &dataValueArray[i], o.engine, jsOwner)
		}
		// TODO Check the bool result and return an error.
		C.objectInvoke(o.addr, cmethod, &result, &dataValueArray[0], C.int(len(params)))
	})
	return unpackDataValue(&result, o.engine)
}

// Destroy finalizes the value and releases any resources used.
// The value must not be used after calling this method.
func (o *commonObject) Destroy() {
	// TODO Must protect against destroyment when object isn't owned.
	gui(func() {
		if o.addr != nilPtr {
			C.delObject(o.addr)
		}
	})
}

// Window represents a QML window where components are rendered.
type Window struct {
	commonObject
}

// Show exposes the window.
func (w *Window) Show() {
	gui(func() {
		C.viewShow(w.addr)
	})
}

// Hide hides the window.
func (w *Window) Hide() {
	gui(func() {
		C.viewHide(w.addr)
	})
}

// Root returns the root component instance being rendered in the window.
func (w *Window) Root() *Value {
	// XXX Test this.
	var object Value
	gui(func() {
		object.addr = C.viewRootObject(w.addr)
	})
	return &object
}

// Wait blocks the current goroutine until the window is closed.
func (w *Window) Wait() {
	// XXX Test this.
	var m sync.Mutex
	m.Lock()
	gui(func() {
		// TODO Must be able to wait for the same Window from multiple goroutines.
		// type foo { m sync.Mutex; next *foo }
		// TODO If the window is not visible, must return immediately.
		waitingWindows[w.addr] = &m
		C.viewConnectHidden(w.addr)
	})
	m.Lock()
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

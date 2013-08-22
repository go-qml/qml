package qml

// #include <stdlib.h>
//
// #include "capi.h"
//
import "C"

import (
	"errors"
	"strings"
	"unsafe"
)

// InitOptions holds options to initialize the qml package.
type InitOptions struct {
	// Reserved for coming options.
}

var initialized = false

// Init initializes the qml package with the provided parameters.
// If the options parameter is nil, default options suitable for a
// normal graphic application will be used.
//
// Init must be called only once, and before any other QML functionality is used.
func Init(options *InitOptions) {
	if initialized {
		panic("qml.Init called more than once")
	}
	initialized = true

	go guiLoop()

	// Wait for app to be created and event loop to be running.
	gui(func() {})
}

// Engine provides an environment for instantiating QML components.
type Engine struct {
	addr unsafe.Pointer
}

func (e *Engine) assertValid() {
	if e.addr == nilPtr {
		panic("engine already closed")
	}
}

// NewEngine returns a new QML engine.
//
// The Close method must be called to release the resources
// used by the engine when done using it.
func NewEngine() *Engine {
	var engine Engine
	gui(func() {
		engine.addr = C.newEngine(nil)
	})
	return &engine
}

// Close releases resources used by the engine. The engine must
// not be used after calling it.
func (e *Engine) Close() {
	if e.addr != nilPtr {
		gui(func() {
			C.delEngine(e.addr)
		})
		e.addr = nilPtr;
	}
}

func (e *Engine) RootContext() *Context {
	e.assertValid()
	var context Context
	gui(func() {
		context.addr = C.engineRootContext(e.addr)
	})
	return &context
}

type Context struct {
	addr unsafe.Pointer
}

type reference struct {
	ifacep *interface{}
	valuep unsafe.Pointer
}

var refs = make(map[interface{}]reference)

func (c *Context) Set(name string, value interface{}) {
	cname, cnamelen := unsafeStringData(name)
	gui(func() {
		var dvalue C.DataValue
		packDataValue(value, &dvalue)

		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)

		C.contextSetProperty(c.addr, qname, &dvalue)
	})
}

func (c *Context) SetObject(value interface{}) {
	// TODO This is leaking. Must figure how to decref the QObject when the context is done with it,
	// so that we can decref it locally as well, and drop the map reference when it reaches zero.
	// Must also lock refs.
	gui(func() {
		ref, ok := refs[value]
		if !ok {
			ref.ifacep = &value
			ref.valuep = C.newValue(unsafe.Pointer(&value), typeInfo(value))
			refs[value] = ref
		}
		C.contextSetObject(c.addr, ref.valuep)
	})
}

func (c *Context) Get(name string) interface{} {
	cname, cnamelen := unsafeStringData(name)

	var dvalue C.DataValue
	gui(func() {
		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)

		C.contextGetProperty(c.addr, qname, &dvalue)
	})
	return unpackDataValue(&dvalue)
}

type Component struct {
	addr unsafe.Pointer
}

// TODO What's a nice way to delete the component and created component objects?

func NewComponent(engine *Engine) *Component {
	var component Component
	gui(func() {
		component.addr = C.newComponent(engine.addr, nilPtr)
	})
	return &component
}

func (c *Component) SetData(location string, data []byte) error {
	cdata, cdatalen := unsafeBytesData(data)
	cloc, cloclen := unsafeStringData(location)
	var err error
	gui(func() {
		C.componentSetData(c.addr, cdata, cdatalen, cloc, cloclen)
		message := C.componentErrorString(c.addr)
		if message != nilCharPtr {
			err = errors.New(strings.TrimRight(C.GoString(message), "\n"))
			C.free(unsafe.Pointer(message))
		}
	})
	return err
}

func (c *Component) Create(context *Context) *Object {
	// TODO Destroy object.
	var object Object
	gui(func() {
		object.addr = C.componentCreate(c.addr, context.addr)
	})
	return &object
}

func (c *Component) CreateWindow(context *Context) *Window {
	var window Window
	gui(func() {
		window.addr = C.componentCreateView(c.addr, context.addr)
	})
	return &window
}

type Object struct {
	addr unsafe.Pointer
}

func (o *Object) Get(property string) interface{} {
	cproperty := C.CString(property)
	defer C.free(unsafe.Pointer(cproperty))

	var value C.DataValue
	gui(func() {
		C.objectGetProperty(o.addr, cproperty, &value)
	})
	return unpackDataValue(&value)
}

//func NewWindow(engine *qml.Engine) {
//	return &Window{C.newView(engine)}
//}

type Window struct {
	addr unsafe.Pointer
}

func (w *Window) Show() {
	gui(func() {
		C.viewShow(w.addr)
	})
}

package qml

// #cgo CPPFLAGS: -I/usr/include/qt5/QtCore/5.0.2/QtCore
// #cgo pkg-config: Qt5Core Qt5Widgets Qt5Quick
// #cgo LDFLAGS: -lstdc++
//
// #include <stdlib.h>
//
// #include "capi.h"
//
import "C"

import (
	"errors"
	"reflect"
	"strings"
	"unsafe"
)

// InitOptions holds options to initialize the qml package.
type InitOptions struct {
	// Reserved for coming options.
}

var initialized = false

var qapp unsafe.Pointer

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

	// Must not be de-allocated according to QApp's docs.
	argv := (**C.char)(C.malloc(ptrSize * 2))
	*(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + uintptr(ptrSize * 0))) = C.CString("")
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + uintptr(ptrSize * 1))) = 0

	qapp = C.newGuiApplication(1, argv)
}

// Run runs the main QML event loop.
func Run() {
	C.applicationExec(qapp);
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
	return &Engine{C.newEngine(nil)}
}

// Close releases resources used by the engine. The engine must
// not be used after calling it.
func (e *Engine) Close() {
	if e.addr != nilPtr {
		C.delEngine(e.addr)
		e.addr = nilPtr;
	}
}

func (e *Engine) RootContext() *Context {
	e.assertValid()
	return &Context{C.engineRootContext(e.addr)}
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
	qname := C.newString(cname, cnamelen)
	defer C.delString(qname)

	dvalue := packDataValue(value)
	C.contextSetProperty(c.addr, qname, dvalue)
}

func (c *Context) SetObject(value interface{}) {
	// TODO This is leaking. Must figure how to decref the QObject when the context is done with it,
	// so that we can decref it locally as well, and drop the map reference when it reaches zero.
	// Must also lock refs.
	ref, ok := refs[value]
	if !ok {
		ref.ifacep = &value
		ref.valuep = C.newValue(unsafe.Pointer(&value), typeInfo(value))
		refs[value] = ref
	}
	C.contextSetObject(c.addr, ref.valuep)
}

func (c *Context) Get(name string) interface{} {
	cname, cnamelen := unsafeStringData(name)
	qname := C.newString(cname, cnamelen)
	defer C.delString(qname)

	var dvalue C.DataValue
	C.contextGetProperty(c.addr, qname, &dvalue)
	return unpackDataValue(&dvalue)
}

type Component struct {
	addr unsafe.Pointer
}

func NewComponent(engine *Engine) *Component {
	return &Component{C.newComponent(engine.addr, nilPtr)}
}

func (c *Component) SetData(location string, data []byte) error {
	cdata, cdatalen := unsafeBytesData(data)
	cloc, cloclen := unsafeStringData(location)
	C.componentSetData(c.addr, cdata, cdatalen, cloc, cloclen)
	message := C.componentErrorString(c.addr)
	if message != nilCharPtr {
		err := errors.New(strings.TrimRight(C.GoString(message), "\n"))
		C.free(unsafe.Pointer(message))
		return err
	}
	return nil
}

func (c *Component) Create(context *Context) *Object {
	// TODO Destroy object.
	return &Object{C.componentCreate(c.addr, context.addr)}
}

type Object struct {
	addr unsafe.Pointer
}

func (o *Object) Get(property string) interface{} {
	cproperty := C.CString(property)
	defer C.free(unsafe.Pointer(cproperty))

	var value C.DataValue
	C.objectGetProperty(o.addr, cproperty, &value)
	return unpackDataValue(&value)
}

// TODO What's a nice way to delete the component and created component objects?

//export hookReadField
func hookReadField(ptr unsafe.Pointer, memberIndex C.int, result unsafe.Pointer) {
	ifacep := (*interface{})(ptr)
	//fmt.Printf("QML requested member %d for Go's %T at %p.\n", memberIndex, *ifacep, ifacep)
	field := reflect.ValueOf(*ifacep).Elem().Field(int(memberIndex))

	switch field.Type().Kind() {
	case reflect.String:
		*(**C.char)(result) = C.CString(field.String()) // XXX This is leaking.
	case reflect.Bool:
		var b int32
		if field.Bool() {
			b = 1
		}
		*(*int32)(result) = b
	case reflect.Int:
		if !intIs64 {
			*(*int32)(result) = int32(field.Int())
			break
		}
		fallthrough
	case reflect.Int64:
		*(*int64)(result) = field.Int()
	case reflect.Int32:
		*(*int32)(result) = int32(field.Int())
	case reflect.Float64:
		*(*float64)(result) = field.Float()
	case reflect.Float32:
		*(*float32)(result) = float32(field.Float())
	default:
		panic("gqReadField got unsupported type: " + field.Type().Name())
	}
}

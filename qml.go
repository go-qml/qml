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
	"fmt"
	"reflect"
	"unsafe"
)

// InitOptions holds options to initialize the qml package.
type InitOptions struct {
	// Reserved for coming options.
}

var initialized = false

var qapp unsafe.Pointer

var wordSize = C.size_t(unsafe.Sizeof(uintptr(0)))

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
	argv := (**C.char)(C.malloc(wordSize * 2))
	*(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + uintptr(wordSize * 0))) = C.CString("")
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + uintptr(wordSize * 1))) = 0

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
	if e.addr == invalidPointer {
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

var invalidPointer = unsafe.Pointer(uintptr(0));

// Close releases resources used by the engine. The engine must
// not be used after calling it.
func (e *Engine) Close() {
	if e.addr != invalidPointer {
		C.delEngine(e.addr)
		e.addr = invalidPointer;
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

func newString(s string) unsafe.Pointer {
	// TODO Test the s == "" case.
	return C.newString(*(**C.char)(unsafe.Pointer(&s)), C.int(len(s)))
}

func (c *Context) Set(name string, value interface{}) {
	qname := newString(name)
	defer C.delString(qname)

	switch value := value.(type) {
	case string:
		qvalue := newString(value)
		C.contextSetPropertyString(c.addr, qname, qvalue)
		C.delString(qvalue)
		return
	case bool:
		var b C.int32_t
		if value {
			b = 1
		}
		C.contextSetPropertyBool(c.addr, qname, b)
		return
	case int64:
		C.contextSetPropertyInt64(c.addr, qname, C.int64_t(value))
		return
	case int32:
		C.contextSetPropertyInt32(c.addr, qname, C.int32_t(value))
		return
	case float64:
		C.contextSetPropertyFloat64(c.addr, qname, C.double(value))
		return
	case float32:
		C.contextSetPropertyFloat32(c.addr, qname, C.float(value))
		return
	}

	// TODO This is leaking. Must figure how to decref the QObject when the context is done with it,
	// so that we can decref it locally as well, and drop the map reference when it reaches zero.
	// Must also lock refs.
	ref, ok := refs[value]
	if !ok {
		ref.ifacep = &value
		ref.valuep = C.newValue(unsafe.Pointer(&value), typeInfo(value))
		refs[value] = ref
	}
	C.contextSetPropertyObject(c.addr, qname, ref.valuep)
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
	qname := newString(name)
	defer C.delString(qname)

	var mem int64
	var result = unsafe.Pointer(&mem)
	var dtype C.DataType
	C.contextGetProperty(c.addr, qname, result, &dtype)

	switch dtype {
	case C.DTString:
		s := C.GoString(*(**C.char)(result))
		C.free(unsafe.Pointer(*(**C.char)(result)))
		return s
	case C.DTBool:
		if *(*int32)(result) == 0 {
			return false
		}
		return true
	case C.DTInt64:
		return *(*int64)(result)
	case C.DTInt32:
		return *(*int32)(result)
	case C.DTFloat64:
		return *(*float64)(result)
	case C.DTFloat32:
		return *(*float32)(result)
	case C.DTGoAddr:
		return **(**interface{})(result)
	}

	panic(fmt.Sprintf("unsupported data type: %d", dtype))
}


//export hookReadField
func hookReadField(ptr unsafe.Pointer, memberIndex C.int, result unsafe.Pointer) {
	ifacep := (*interface{})(ptr)
	fmt.Printf("QML requested member %d for Go's %T at %p.\n", memberIndex, *ifacep, ifacep)
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

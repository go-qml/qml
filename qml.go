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
	"unicode"
	"unsafe"
)

// InitOptions holds options to initialize the qml package.
type InitOptions struct {
	// Reserved for coming options.
}

var initialized = false

var gqApp unsafe.Pointer

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

	gqApp = C.newGuiApplication(1, argv)
}

// Run runs the main QML event loop.
func Run() {
	C.applicationExec(gqApp);
}

// Engine provides an environment for instantiating QML components.
type Engine struct {
	addr unsafe.Pointer
}

// NewEngine returns a new QML engine.
func NewEngine() *Engine {
	return &Engine{C.newEngine(nil)}
}

type Context struct {
	addr unsafe.Pointer
}

func (e *Engine) RootContext() *Context {
	return &Context{C.engineRootContext(e.addr)}
}

type reference struct {
	goValue *interface{}
	gqValue unsafe.Pointer
}

var refs = make(map[interface{}]reference)

func newString(s string) unsafe.Pointer {
	return C.newString(*(**C.char)(unsafe.Pointer(&s)), C.int(len(s)))
}

func (c *Context) Set(name string, value interface{}) {
	// TODO Must handle the name == "" case.
	qname := newString(name)
	defer C.delString(qname)

	switch value := value.(type) {
	case string:
		// TODO Must handle the value == "" case.
		qvalue := C.newString(*(**C.char)(unsafe.Pointer(&value)), C.int(len(value)))
		C.contextSetString(c.addr, qname, qvalue)
		C.delString(qvalue)
		return
	case int64:
		C.contextSetInt64(c.addr, qname, C.int64_t(value))
		return
	case int32:
		C.contextSetInt32(c.addr, qname, C.int32_t(value))
		return
	default:
		panic(fmt.Sprintf("Context.Set of %T is unsupported for now", value))
	}

	// TODO This is leaking. Must figure how to decref the QObject when the context is done with it,
	// so that we can decref it locally as well, and drop the map reference when it reaches zero.
	// Must also lock refs.
	ref, ok := refs[value]
	if !ok {
		ref.goValue = &value
		ref.gqValue = C.newValue(unsafe.Pointer(&value), typeInfo(value))
		refs[value] = ref
	}
	C.contextSetObject(c.addr, qname, ref.gqValue)
}

func (c *Context) Get(name string) interface{} {
	// Do this in a more efficient way.
	qname := newString(name)
	defer C.delString(qname)

	var mem int64
	var result = unsafe.Pointer(&mem)
	var dtype C.DataType
	C.contextGet(c.addr, qname, result, &dtype)

	switch dtype {
	case C.DTString:
		s := C.GoString(*(**C.char)(result))
		C.free(unsafe.Pointer(*(**C.char)(result)))
		return s
	case C.DTInt64:
		return *(*int64)(result)
	case C.DTInt32:
		return *(*int32)(result)
	case C.DTFloat64:
		return *(*float64)(result)
	case C.DTFloat32:
		return *(*float32)(result)
	}

	panic(fmt.Sprintf("unsupported data type: %d", dtype))
}


var typeInfoSize = C.size_t(unsafe.Sizeof(C.GoTypeInfo{}))
var memberInfoSize = C.size_t(unsafe.Sizeof(C.GoMemberInfo{}))

func typeInfo(v interface{}) *C.GoTypeInfo {
	vt := reflect.TypeOf(v)
	for vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}

	typeInfo := (*C.GoTypeInfo)(C.malloc(typeInfoSize))
	typeInfo.typeName = C.CString(vt.Name())

	numField := vt.NumField()

	// struct { FooBar T; Baz T } => "fooBar\0baz\0"
	namesLen := 0
	for i := 0; i < numField; i++ {
		namesLen += len(vt.Field(i).Name)
	}
	names := make([]byte, 0, namesLen)
	for i := 0; i < numField; i++ {
		name := vt.Field(i).Name
		for i, rune := range name {
			if i == 0 {
				names = append(names, string(unicode.ToLower(rune))...)
			} else {
				names = append(names, name[i:]...)
				break
			}
		}
		names = append(names, 0)
	}
	typeInfo.memberNames = C.CString(string(names))

	// Assemble information on members.
	membersi := uintptr(0)
	mnamesi := uintptr(0)
	members := uintptr(C.malloc(memberInfoSize*C.size_t(numField) + 1))
	mnames := uintptr(unsafe.Pointer(typeInfo.memberNames))
	for i := 0; i < numField; i++ {
		field := vt.Field(i)
		memberInfo := (*C.GoMemberInfo)(unsafe.Pointer(members + (uintptr(memberInfoSize) * membersi)))
		memberInfo.memberName = (*C.char)(unsafe.Pointer(mnames + mnamesi))
		memberInfo.memberType = gqKindFor(field.Type)
		memberInfo.memberIndex = C.int(i)
		membersi += 1
		mnamesi += uintptr(len(field.Name)) + 1
	}
	typeInfo.members = (*C.GoMemberInfo)(unsafe.Pointer(members))
	typeInfo.membersLen = C.int(membersi)
	return typeInfo
}

var (
	intIs64   bool
	intT C.DataType
)

func init() {
	intIs64 = unsafe.Sizeof(int64(0)) == unsafe.Sizeof(int(0))
	if intIs64 {
		intT = C.DTInt64
	} else {
		intT = C.DTInt32
	}
}

func gqKindFor(typ reflect.Type) C.DataType {
	switch typ.Kind() {
	case reflect.String:
		return C.DTString
	case reflect.Int64:
		return C.DTInt64
	case reflect.Int32:
		return C.DTInt32
	case reflect.Int:
		return intT
	case reflect.Float32:
		return C.DTFloat32
	case reflect.Float64:
		return C.DTFloat64
	}
	panic("Go type not supported yet: " + typ.Name())
}

//export hookReadField
func hookReadField(ptr unsafe.Pointer, memberIndex C.int, result unsafe.Pointer) {
	ifacep := (*interface{})(ptr)
	fmt.Printf("QML requested member %d for Go's %T at %p.\n", memberIndex, *ifacep, ifacep)
	field := reflect.ValueOf(*ifacep).Field(int(memberIndex))

	switch field.Type().Kind() {
	case reflect.String:
		*(**C.char)(result) = C.CString(field.String()) // XXX This is leaking.
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

package qml

// #include <stdlib.h>
// #include "capi.h"
import "C"

import (
	"fmt"
	"reflect"
	"unicode"
	"unsafe"
)

var (
	intIs64 bool
	intDT   C.DataType

	ptrSize = C.size_t(unsafe.Sizeof(uintptr(0)))

	nilPtr = unsafe.Pointer(uintptr(0))
	nilCharPtr = (*C.char)(nilPtr)

	typeString = reflect.TypeOf("")
	typeBool = reflect.TypeOf(false)
	typeInt = reflect.TypeOf(int(0))
	typeInt64 = reflect.TypeOf(int64(0))
	typeInt32 = reflect.TypeOf(int32(0))
	typeFloat64 = reflect.TypeOf(float64(0))
	typeFloat32 = reflect.TypeOf(float32(0))
)

func init() {
	var i int = 1<<31 - 1
	intIs64 = (i+1 > 0)
	if intIs64 {
		intDT = C.DTInt64
	} else {
		intDT = C.DTInt32
	}
}

func dataTypeOf(typ reflect.Type) C.DataType {
	switch typ.Kind() {
	case reflect.String:
		return C.DTString
	case reflect.Bool:
		return C.DTBool
	case reflect.Int:
		return intDT
	case reflect.Int64:
		return C.DTInt64
	case reflect.Int32:
		return C.DTInt32
	case reflect.Float32:
		return C.DTFloat32
	case reflect.Float64:
		return C.DTFloat64
	}
	panic("Go type not supported yet: " + typ.Name())
}

func packDataValue(value interface{}) *C.DataValue {
	var dvalue C.DataValue
	datap := unsafe.Pointer(&dvalue.data)
	switch value := value.(type) {
	case string:
		dvalue.dataType = C.DTString
		cstr, cstrlen := unsafeStringData(value)
		*(**C.char)(datap) = cstr
		dvalue.len = cstrlen
	case bool:
		dvalue.dataType = C.DTBool
		*(*bool)(datap) = value
	case int:
		dvalue.dataType = intDT
		*(*int)(datap) = value
	case int64:
		dvalue.dataType = C.DTInt64
		*(*int64)(datap) = value
	case int32:
		dvalue.dataType = C.DTInt32
		*(*int32)(datap) = value
	case float64:
		dvalue.dataType = C.DTFloat64
		*(*float64)(datap) = value
	case float32:
		dvalue.dataType = C.DTFloat32
		*(*float32)(datap) = value
	default:
		// TODO This is leaking. Must figure how to decref the QObject when the context is done with it,
		// so that we can decref it locally as well, and drop the map reference when it reaches zero.
		// Must also lock refs.
		ref, ok := refs[value]
		if !ok {
			var value interface{} = value
			ref.ifacep = &value
			ref.valuep = C.newValue(unsafe.Pointer(ref.ifacep), typeInfo(value))
			refs[value] = ref
		}
		dvalue.dataType = C.DTObject
		*(*unsafe.Pointer)(datap) = ref.valuep
	}
	return &dvalue
}

func unpackDataValue(dvalue *C.DataValue) interface{} {
	datap := unsafe.Pointer(&dvalue.data)
	switch dvalue.dataType {
	case C.DTString:
		s := C.GoStringN(*(**C.char)(datap), dvalue.len)
		C.free(unsafe.Pointer(*(**C.char)(datap)))
		return s
	case C.DTBool:
		return *(*bool)(datap)
	case C.DTInt64:
		return *(*int64)(datap)
	case C.DTInt32:
		return *(*int32)(datap)
	case C.DTFloat64:
		return *(*float64)(datap)
	case C.DTFloat32:
		return *(*float32)(datap)
	case C.DTGoAddr:
		return **(**interface{})(datap)
	}
	panic(fmt.Sprintf("unsupported data type: %d", dvalue.dataType))
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
		memberInfo.memberType = dataTypeOf(field.Type)
		memberInfo.memberIndex = C.int(i)
		membersi += 1
		mnamesi += uintptr(len(field.Name)) + 1
	}
	typeInfo.members = (*C.GoMemberInfo)(unsafe.Pointer(members))
	typeInfo.membersLen = C.int(membersi)
	return typeInfo
}

func unsafeBytesData(b []byte) (*C.char, C.int) {
	return *(**C.char)(unsafe.Pointer(&b)), C.int(len(b))
}

func unsafeStringData(s string) (*C.char, C.int) {
	return *(**C.char)(unsafe.Pointer(&s)), C.int(len(s))
}

// unsafeString returns a Go string backed by C data.
//
// The returned string must be used only in the implementation
// of the qml package, since its data cannot be relied upon.
func unsafeString(data *C.char, size C.int) string {
	var s string
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh.Data = uintptr(unsafe.Pointer(data))
	sh.Len = int(size)
	return s
}

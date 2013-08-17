package qml

// #include <stdlib.h>
// #include "capi.h"
import "C"

import (
	"reflect"
	"unicode"
	"unsafe"
)


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

func dataTypeOf(typ reflect.Type) C.DataType {
	switch typ.Kind() {
	case reflect.String:
		return C.DTString
	case reflect.Bool:
		return C.DTBool
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

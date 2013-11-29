package qml

// #include <stdlib.h>
// #include "capi.h"
import "C"

import (
	"bytes"
	"fmt"
	"image/color"
	"reflect"
	"strings"
	"unicode"
	"unsafe"
)

var (
	intIs64 bool
	intDT   C.DataType

	ptrSize = C.size_t(unsafe.Sizeof(uintptr(0)))

	nilPtr     = unsafe.Pointer(uintptr(0))
	nilCharPtr = (*C.char)(nilPtr)

	typeString   = reflect.TypeOf("")
	typeBool     = reflect.TypeOf(false)
	typeInt      = reflect.TypeOf(int(0))
	typeInt64    = reflect.TypeOf(int64(0))
	typeInt32    = reflect.TypeOf(int32(0))
	typeFloat64  = reflect.TypeOf(float64(0))
	typeFloat32  = reflect.TypeOf(float32(0))
	typeIface    = reflect.TypeOf(new(interface{})).Elem()
	typeRGBA     = reflect.TypeOf(color.RGBA{})
	typeObjSlice = reflect.TypeOf([]Object(nil))
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

// packDataValue packs the provided Go value into a C.DataValue for
// shiping into C++ land.
//
// For simple types (bool, int, etc) value is converted into a
// native C++ value. For anything else, including cases when value
// has a type that has an underlying simple type, the Go value itself
// is encapsulated into a C++ wrapper so that field access and method
// calls work.
//
// This must be run from the main GUI thread due to the cases where
// calling wrapGoValue is necessary.
func packDataValue(value interface{}, dvalue *C.DataValue, engine *Engine, owner valueOwner) {
	datap := unsafe.Pointer(&dvalue.data)
	if value == nil {
		dvalue.dataType = C.DTInvalid
		return
	}
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
		if value > 1<<31-1 {
			dvalue.dataType = C.DTInt64
			*(*int64)(datap) = int64(value)
		} else {
			dvalue.dataType = C.DTInt32
			*(*int32)(datap) = int32(value)
		}
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
	case *Common:
		dvalue.dataType = C.DTObject
		*(*unsafe.Pointer)(datap) = value.addr
	case color.RGBA:
		dvalue.dataType = C.DTColor
		*(*uint32)(datap) = uint32(value.A)<<24 | uint32(value.R)<<16 | uint32(value.G)<<8 | uint32(value.B)
	default:
		dvalue.dataType = C.DTObject
		if obj, ok := value.(Object); ok {
			*(*unsafe.Pointer)(datap) = obj.Common().addr
		} else {
			*(*unsafe.Pointer)(datap) = wrapGoValue(engine, value, owner)
		}
	}
}

// TODO Handle byte slices.

// unpackDataValue converts a value shipped by C++ into a native Go value.
//
// HEADS UP: This is considered safe to be run out of the main GUI thread.
//           If that changes, fix the call sites.
func unpackDataValue(dvalue *C.DataValue, engine *Engine) interface{} {
	datap := unsafe.Pointer(&dvalue.data)
	switch dvalue.dataType {
	case C.DTString:
		s := C.GoStringN(*(**C.char)(datap), dvalue.len)
		// TODO If we move all unpackDataValue calls to the GUI thread,
		// can we get rid of this allocation somehow?
		C.free(unsafe.Pointer(*(**C.char)(datap)))
		return s
	case C.DTBool:
		return *(*bool)(datap)
	case C.DTInt64:
		return *(*int64)(datap)
	case C.DTInt32:
		return int(*(*int32)(datap))
	case C.DTFloat64:
		return *(*float64)(datap)
	case C.DTFloat32:
		return *(*float32)(datap)
	case C.DTColor:
		var c uint32 = *(*uint32)(datap)
		return color.RGBA{byte(c >> 16), byte(c >> 8), byte(c), byte(c >> 24)}
	case C.DTGoAddr:
		return (*(**valueFold)(datap)).gvalue
	case C.DTInvalid:
		return nil
	case C.DTObject:
		// TODO Would be good to preserve identity on the Go side.
		return &Common{
			engine: engine,
			addr:   *(*unsafe.Pointer)(datap),
		}
	case C.DTValueList:
		var dvlist []C.DataValue
		var dvlisth = (*reflect.SliceHeader)(unsafe.Pointer(&dvlist))
		dvlisth.Data = uintptr(*(*unsafe.Pointer)(datap))
		dvlisth.Len = int(dvalue.len)
		dvlisth.Cap = int(dvalue.len)
		result := make([]interface{}, len(dvlist))
		for i := range result {
			result[i] = unpackDataValue(&dvlist[i], engine)
		}
		C.free(*(*unsafe.Pointer)(datap))
		return &List{result}
	}
	panic(fmt.Sprintf("unsupported data type: %d", dvalue.dataType))
}

func dataTypeOf(typ reflect.Type) C.DataType {
	// Compare against the specific types rather than their kind.
	// Custom types may have methods that must be supported.
	switch typ {
	case typeString:
		return C.DTString
	case typeBool:
		return C.DTBool
	case typeInt:
		return intDT
	case typeInt64:
		return C.DTInt64
	case typeInt32:
		return C.DTInt32
	case typeFloat32:
		return C.DTFloat32
	case typeFloat64:
		return C.DTFloat64
	case typeIface:
		return C.DTAny
	case typeRGBA:
		return C.DTColor
	case typeObjSlice:
		return C.DTListProperty
	}
	return C.DTObject
}

var typeInfoSize = C.size_t(unsafe.Sizeof(C.GoTypeInfo{}))
var memberInfoSize = C.size_t(unsafe.Sizeof(C.GoMemberInfo{}))

var typeInfoCache = make(map[reflect.Type]*C.GoTypeInfo)

func typeInfo(v interface{}) *C.GoTypeInfo {
	vt := reflect.TypeOf(v)
	for vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}

	typeInfo := typeInfoCache[vt]
	if typeInfo != nil {
		return typeInfo
	}

	typeInfo = (*C.GoTypeInfo)(C.malloc(typeInfoSize))
	typeInfo.typeName = C.CString(vt.Name())
	typeInfo.metaObject = nilPtr

	var onChanged map[string]int

	// TODO Only do that if it's a struct?
	vtptr := reflect.PtrTo(vt)

	numField := vt.NumField()
	prvField := 0
	numMethod := vtptr.NumMethod()

	// struct { FooBar T; Baz T } => "fooBar\0baz\0"
	namesLen := 0
	for i := 0; i < numField; i++ {
		field := vt.Field(i)
		if field.PkgPath != "" {
			prvField++ // not exported
			continue
		}
		namesLen += len(field.Name) + 1
	}
	for i := 0; i < numMethod; i++ {
		namesLen += len(vtptr.Method(i).Name) + 1
	}
	names := make([]byte, 0, namesLen)
	for i := 0; i < numField; i++ {
		field := vt.Field(i)
		if field.PkgPath != "" {
			continue // not exported
		}
		name := field.Name
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
	for i := 0; i < numMethod; i++ {
		name := vtptr.Method(i).Name
		for i, rune := range name {
			if i == 0 {
				names = append(names, string(unicode.ToLower(rune))...)
			} else {
				names = append(names, name[i:]...)
				break
			}
		}
		names = append(names, 0)

		// Track "On*Changed" notification methods.
		if len(name) > 9 && name[0] == 'O' && name[1] == 'n' && strings.HasSuffix(name, "Changed") {
			if onChanged == nil {
				onChanged = make(map[string]int)
			}
			onChanged[name[2:len(name)-7]] = i
		}
	}
	if len(names) != namesLen {
		panic("pre-allocated buffer size was wrong")
	}
	typeInfo.memberNames = C.CString(string(names))

	// Assemble information on members.
	membersLen := numField - prvField + numMethod
	membersi := uintptr(0)
	mnamesi := uintptr(0)
	members := uintptr(C.malloc(memberInfoSize * C.size_t(membersLen)))
	mnames := uintptr(unsafe.Pointer(typeInfo.memberNames))
	for i := 0; i < numField; i++ {
		field := vt.Field(i)
		if field.PkgPath != "" {
			continue // not exported
		}
		memberInfo := (*C.GoMemberInfo)(unsafe.Pointer(members + uintptr(memberInfoSize)*membersi))
		memberInfo.memberName = (*C.char)(unsafe.Pointer(mnames + mnamesi))
		memberInfo.memberType = dataTypeOf(field.Type)
		memberInfo.reflectIndex = C.int(i)
		memberInfo.reflectChangedIndex = -1
		memberInfo.addrOffset = C.int(field.Offset)
		membersi += 1
		mnamesi += uintptr(len(field.Name)) + 1
		if methodIndex, ok := onChanged[field.Name]; ok {
			memberInfo.reflectChangedIndex = C.int(methodIndex)
		}
	}
	for i := 0; i < numMethod; i++ {
		method := vtptr.Method(i)
		memberInfo := (*C.GoMemberInfo)(unsafe.Pointer(members + uintptr(memberInfoSize)*membersi))
		memberInfo.memberName = (*C.char)(unsafe.Pointer(mnames + mnamesi))
		memberInfo.memberType = C.DTMethod
		memberInfo.reflectIndex = C.int(i)
		memberInfo.reflectChangedIndex = -1
		memberInfo.addrOffset = 0
		signature, result := methodQtSignature(method)
		// TODO The signature data might be embedded in the same array as the member names.
		memberInfo.methodSignature = C.CString(signature)
		memberInfo.resultSignature = C.CString(result)
		// TODO Sort out methods with a variable number of arguments.
		// It's called while bound, so drop the receiver.
		memberInfo.numIn = C.int(method.Type.NumIn() - 1)
		memberInfo.numOut = C.int(method.Type.NumOut())
		membersi += 1
		mnamesi += uintptr(len(method.Name)) + 1
	}
	typeInfo.members = (*C.GoMemberInfo)(unsafe.Pointer(members))
	typeInfo.membersLen = C.int(membersLen)

	typeInfo.fields = typeInfo.members
	typeInfo.fieldsLen = C.int(numField - prvField)
	typeInfo.methods = (*C.GoMemberInfo)(unsafe.Pointer(members + uintptr(memberInfoSize)*uintptr(typeInfo.fieldsLen)))
	typeInfo.methodsLen = C.int(numMethod)

	if int(membersi) != membersLen {
		panic("used more space than allocated for member names")
	}
	if int(mnamesi) != namesLen {
		panic("allocated buffer doesn't match used space")
	}
	if typeInfo.fieldsLen+typeInfo.methodsLen != typeInfo.membersLen {
		panic("lengths are inconsistent")
	}

	typeInfoCache[vt] = typeInfo
	return typeInfo
}

func methodQtSignature(method reflect.Method) (signature, result string) {
	var buf bytes.Buffer
	for i, rune := range method.Name {
		if i == 0 {
			buf.WriteRune(unicode.ToLower(rune))
		} else {
			buf.WriteString(method.Name[i:])
			break
		}
	}
	buf.WriteByte('(')
	n := method.Type.NumIn()
	for i := 1; i < n; i++ {
		if i > 1 {
			buf.WriteByte(',')
		}
		buf.WriteString("QVariant")
	}
	buf.WriteByte(')')
	signature = buf.String()

	switch method.Type.NumOut() {
	case 0:
		// keep it as ""
	case 1:
		result = "QVariant"
	default:
		result = "QVariantList"
	}
	return
}

func hashable(value interface{}) (hashable bool) {
	defer func() { recover() }()
	return value == value
}

// unsafeString returns a Go string backed by C data.
//
// If the C data is deallocated or moved, the string will be
// invalid and will crash the program if used. As such, the
// resulting string must only be used inside the implementation
// of the qml package and while the life time of the C data
// is guaranteed.
func unsafeString(data *C.char, size C.int) string {
	var s string
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh.Data = uintptr(unsafe.Pointer(data))
	sh.Len = int(size)
	return s
}

// unsafeStringData returns a C string backed by Go data. The C
// string is NOT null-terminated, so its length must be taken
// into account.
//
// If the s Go string is garbage collected, the returned C data
// will be invalid and will crash the program if used. As such,
// the resulting data must only be used inside the implementation
// of the qml package and while the life time of the Go string
// is guaranteed.
func unsafeStringData(s string) (*C.char, C.int) {
	return *(**C.char)(unsafe.Pointer(&s)), C.int(len(s))
}

// unsafeBytesData returns a C string backed by Go data. The C
// string is NOT null-terminated, so its length must be taken
// into account.
//
// If the array backing the b Go slice is garbage collected, the
// returned C data will be invalid and will crash the program if
// used. As such, the resulting data must only be used inside the
// implementation of the qml package and while the life time of
// the Go array is guaranteed.
func unsafeBytesData(b []byte) (*C.char, C.int) {
	return *(**C.char)(unsafe.Pointer(&b)), C.int(len(b))
}

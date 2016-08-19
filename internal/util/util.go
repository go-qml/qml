package util

import (
	"reflect"
	"unsafe"
)

func Hashable(value interface{}) (hashable bool) {
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
func UnsafeString(data unsafe.Pointer, size int) string {
	var s string
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh.Data = uintptr(data)
	sh.Len = size
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
func UnsafeStringData(s string) (unsafe.Pointer, int) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return unsafe.Pointer(sh.Data), sh.Len
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
func UnsafeBytesData(b []byte) (unsafe.Pointer, int) {
	return *(*unsafe.Pointer)(unsafe.Pointer(&b)), len(b)
}

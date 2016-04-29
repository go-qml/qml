package qml

import (
	"fmt"
	"reflect"
)

// List holds a QML list which may be converted to a Go slice of an
// appropriate type via Convert.
//
// In the future this will also be able to hold a reference
// to QML-owned maps, so they can be mutated in place.
type List struct {
	// In the future this will be able to hold a reference to QML-owned
	// lists, so they can be mutated.
	data []interface{}
}

// Len returns the number of elements in the list.
func (l *List) Len() int {
	return len(l.data)
}

// Convert allocates a new slice and copies the list content into it,
// performing type conversions as possible, and then assigns the result
// to the slice pointed to by sliceAddr.
// Convert panics if the list values are not compatible with the
// provided slice.
func (l *List) Convert(sliceAddr interface{}) {
	toPtr := reflect.ValueOf(sliceAddr)
	if toPtr.Kind() != reflect.Ptr || toPtr.Type().Elem().Kind() != reflect.Slice {
		panic(fmt.Sprintf("List.Convert got a sliceAddr parameter that is not a slice address: %#v", sliceAddr))
	}
	err := convertAndSet(toPtr.Elem(), reflect.ValueOf(l), reflect.Value{})
	if err != nil {
		panic(err.Error())
	}
}

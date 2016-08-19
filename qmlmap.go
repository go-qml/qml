package qml

import (
	"fmt"
	"reflect"
)

// Map holds a QML map which may be converted to a Go map of an
// appropriate type via Convert.
//
// In the future this will also be able to hold a reference
// to QML-owned maps, so they can be mutated in place.
type Map struct {
	data []interface{}
}

// Len returns the number of pairs in the map.
func (m *Map) Len() int {
	return len(m.data) / 2
}

// Convert allocates a new map and copies the content of m property to it,
// performing type conversions as possible, and then assigns the result to
// the map pointed to by mapAddr. Map panics if m contains values that
// cannot be converted to the type of the map at mapAddr.
func (m *Map) Convert(mapAddr interface{}) {
	toPtr := reflect.ValueOf(mapAddr)
	if toPtr.Kind() != reflect.Ptr || toPtr.Type().Elem().Kind() != reflect.Map {
		panic(fmt.Sprintf("Map.Convert got a mapAddr parameter that is not a map address: %#v", mapAddr))
	}
	err := convertAndSet(toPtr.Elem(), reflect.ValueOf(m), reflect.Value{})
	if err != nil {
		panic(err.Error())
	}
}

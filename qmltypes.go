package qml

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

// TypeSpec holds the specification of a QML type that is backed by Go logic.
//
// The type specification must be registered with the RegisterTypes function
// before it will be visible to QML code, as in:
//
//     qml.RegisterTypes("GoExtensions", 1, 0, []qml.TypeSpec{{
//		Init: func(p *Person, obj qml.Object) {},
//     }})
//
// See the package documentation for more details.
//
type TypeSpec struct {
	// Init must be set to a function that is called when QML code requests
	// the creation of a new value of this type. The provided function must
	// have the following type:
	//
	//     func(value *CustomType, object qml.Object)
	//
	// Where CustomType is the custom type being registered. The function will
	// be called with a newly created *CustomType and its respective qml.Object.
	Init interface{}

	// Name optionally holds the identifier the type is known as within QML code,
	// when the registered extension module is imported. If not specified, the
	// name of the Go type provided as the first argument of Init is used instead.
	Name string

	// Singleton defines whether a single instance of the type should be used
	// for all accesses, as a singleton value. If true, all properties of the
	// singleton value are directly accessible under the type name.
	Singleton bool

	private struct{} // Force use of fields by name.
}

var types []*TypeSpec

// RegisterTypes registers the provided list of type specifications for use
// by QML code. To access the registered types, they must be imported from the
// provided location and major.minor version numbers.
//
// For example, with a location "GoExtensions", major 4, and minor 2, this statement
// imports all the registered types in the module's namespace:
//
//     import GoExtensions 4.2
//
// See the documentation on QML import statements for details on these:
//
//     http://qt-project.org/doc/qt-5.0/qtqml/qtqml-syntax-imports.html
//
func RegisterTypes(location string, major, minor int, types []TypeSpec) {
	for i := range types {
		err := registerType(location, major, minor, &types[i])
		if err != nil {
			panic(err)
		}
	}
}

func registerType(location string, major, minor int, spec *TypeSpec) error {
	// Copy and hold a reference to the spec data.
	localSpec := *spec

	f := reflect.ValueOf(localSpec.Init)
	ft := f.Type()
	if ft.Kind() != reflect.Func {
		return fmt.Errorf("TypeSpec.Init must be a function, got %#v", localSpec.Init)
	}
	if ft.NumIn() != 2 {
		return fmt.Errorf("TypeSpec.Init's function must accept two arguments: %s", ft)
	}
	firstArg := ft.In(0)
	if firstArg.Kind() != reflect.Ptr || firstArg.Elem().Kind() == reflect.Ptr {
		return fmt.Errorf("TypeSpec.Init's function must take a pointer type as the second argument: %s", ft)
	}
	if ft.In(1) != typeObject {
		return fmt.Errorf("TypeSpec.Init's function must take qml.Object as the second argument: %s", ft)
	}
	customType := typeInfo(reflect.New(firstArg.Elem()).Interface())
	if localSpec.Name == "" {
		localSpec.Name = firstArg.Elem().Name()
		if localSpec.Name == "" {
			panic("cannot determine registered type name; please provide one explicitly")
		}
	}

	var err error
	RunMain(func() {
		cloc := C.CString(location)
		cname := C.CString(localSpec.Name)
		cres := C.int(0)
		if localSpec.Singleton {
			cres = C.registerSingleton(cloc, C.int(major), C.int(minor), cname, customType, unsafe.Pointer(&localSpec))
		} else {
			cres = C.registerType(cloc, C.int(major), C.int(minor), cname, customType, unsafe.Pointer(&localSpec))
		}
		// It doesn't look like it keeps references to these, but it's undocumented and unclear.
		C.free(unsafe.Pointer(cloc))
		C.free(unsafe.Pointer(cname))
		if cres == -1 {
			err = fmt.Errorf("QML engine failed to register type; invalid type location or name?")
		} else {
			types = append(types, &localSpec)
		}
	})

	return err
}

// RegisterConverter registers the convereter function to be called when a
// value with the provided type name is obtained from QML logic. The function
// must return the new value to be used in place of the original value.
func RegisterConverter(typeName string, converter func(engine *Engine, obj Object) interface{}) {
	if converter == nil {
		delete(converters, typeName)
	} else {
		converters[typeName] = converter
	}
}

var converters = make(map[string]func(engine *Engine, obj Object) interface{})

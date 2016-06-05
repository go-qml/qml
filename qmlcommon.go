package qml

// #include <stdlib.h>
//
// #include "capi.h"
//
import "C"

import (
	"fmt"
	"image/color"
	"os"
	"reflect"
	"unsafe"

	"github.com/limetext/qml-go/internal/util"
)

// Common implements the common behavior of all QML objects.
// It implements the Object interface.
type Common struct {
	addr        unsafe.Pointer
	engine      *Engine
	destroyed   bool
	initialized bool
}

var _ Object = (*Common)(nil)

// CommonOf returns the Common QML value for the QObject at addr.
//
// This is meant for extensions that integrate directly with the
// underlying QML logic.
func CommonOf(addr unsafe.Pointer, engine *Engine) *Common {
	c := &Common{nil, engine, false, false}
	c.setAddr(addr)
	return c
}

func (obj *Common) setAddr(addr unsafe.Pointer) {
	if obj.initialized || obj.addr != nil || obj.destroyed {
		panic("Cannot reuse Common!")
	}
	obj.addr = addr
	obj.initialized = true

	if addr != nil {
		obj.On("destroyed", func() { obj.addr = nil; obj.destroyed = true })
	}
}

func (obj *Common) assertInitialized() {
	if !obj.initialized {
		panic("Use of uninitialized object")
	}
	if obj.destroyed {
		panic("Use of destroyed object")
	}
}

// Common returns obj itself.
//
// This provides access to the underlying *Common for types that
// embed it, when these are used via the Object interface.
func (obj *Common) Common() *Common {
	return obj
}

// TypeName returns the underlying type name for the held value.
func (obj *Common) TypeName() string {
	obj.assertInitialized()
	var name string
	RunMain(func() {
		name = C.GoString(C.objectTypeName(obj.addr))
	})
	return name
}

// Addr returns the QML object address.
//
// This is meant for extensions that integrate directly with the
// underlying QML logic.
func (obj *Common) Addr() uintptr {
	obj.assertInitialized()
	return uintptr(obj.addr)
}

// Interface returns the underlying Go value that is being held by
// the object wrapper.
//
// It is a runtime error to call Interface on values that are not
// backed by a Go value.
func (obj *Common) Interface() interface{} {
	obj.assertInitialized()
	var result interface{}
	var cerr *C.error
	RunMain(func() {
		var fold *valueFold
		if cerr = C.objectGoAddr(obj.addr, (*unsafe.Pointer)(unsafe.Pointer(&fold))); cerr == nil {
			result = fold.gvalue
		}
	})
	cmust(cerr)
	return result
}

// Set changes the named object property to the given value.
func (obj *Common) Set(property string, value interface{}) {
	obj.assertInitialized()
	cproperty := C.CString(property)
	defer C.free(unsafe.Pointer(cproperty))
	var cerr *C.error
	RunMain(func() {
		var dvalue C.DataValue
		packDataValue(value, &dvalue, obj.engine, cppOwner)
		cerr = C.objectSetProperty(obj.addr, cproperty, &dvalue)
	})
	cmust(cerr)
}

// Property returns the current value for a property of the object.
// If the property type is known, type-specific methods such as Int
// and String are more convenient to use.
// Property panics if the property does not exist.
func (obj *Common) Property(name string) interface{} {
	obj.assertInitialized()
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var dvalue C.DataValue
	var found C.int
	RunMain(func() {
		found = C.objectGetProperty(obj.addr, cname, &dvalue)
	})
	if found == 0 {
		panic(fmt.Sprintf("object does not have a %q property", name))
	}
	return unpackDataValue(&dvalue, obj.engine)
}

// Int returns the int value of the named property.
// Int panics if the property cannot be represented as an int.
func (obj *Common) Int(property string) int {
	switch value := obj.Property(property).(type) {
	case int64:
		return int(value)
	case int:
		return value
	case uint64:
		return int(value)
	case uint32:
		return int(value)
	case uintptr:
		return int(value)
	case float32:
		return int(value)
	case float64:
		return int(value)
	default:
		panic(fmt.Sprintf("value of property %q cannot be represented as an int: %#v", property, value))
	}
}

// Int64 returns the int64 value of the named property.
// Int64 panics if the property cannot be represented as an int64.
func (obj *Common) Int64(property string) int64 {
	switch value := obj.Property(property).(type) {
	case int64:
		return value
	case int:
		return int64(value)
	case uint64:
		return int64(value)
	case uint32:
		return int64(value)
	case uintptr:
		return int64(value)
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	default:
		panic(fmt.Sprintf("value of property %q cannot be represented as an int64: %#v", property, value))
	}
}

// Float64 returns the float64 value of the named property.
// Float64 panics if the property cannot be represented as float64.
func (obj *Common) Float64(property string) float64 {
	switch value := obj.Property(property).(type) {
	case int64:
		return float64(value)
	case int:
		return float64(value)
	case uint64:
		return float64(value)
	case uint32:
		return float64(value)
	case uintptr:
		return float64(value)
	case float32:
		return float64(value)
	case float64:
		return value
	default:
		panic(fmt.Sprintf("value of property %q cannot be represented as a float64: %#v", property, value))
	}
}

// Bool returns the bool value of the named property.
// Bool panics if the property is not a bool.
func (obj *Common) Bool(property string) bool {
	value := obj.Property(property)
	if b, ok := value.(bool); ok {
		return b
	}
	panic(fmt.Sprintf("value of property %q is not a bool: %#v", property, value))
}

// String returns the string value of the named property.
// String panics if the property is not a string.
func (obj *Common) String(property string) string {
	value := obj.Property(property)
	if s, ok := value.(string); ok {
		return s
	}
	panic(fmt.Sprintf("value of property %q is not a string: %#v", property, value))
}

// Color returns the RGBA value of the named property.
// Color panics if the property is not a color.
func (obj *Common) Color(property string) color.RGBA {
	value := obj.Property(property)
	c, ok := value.(color.RGBA)
	if !ok {
		panic(fmt.Sprintf("value of property %q is not a color: %#v", property, value))
	}
	return c
}

// Object returns the object value of the named property.
// Object panics if the property is not a QML object.
func (obj *Common) Object(property string) Object {
	value := obj.Property(property)
	object, ok := value.(Object)
	if !ok {
		panic(fmt.Sprintf("value of property %q is not a QML object: %#v", property, value))
	}
	return object
}

// List returns the list value of the named property.
// List panics if the property is not a list.
func (obj *Common) List(property string) *List {
	value := obj.Property(property)
	m, ok := value.(*List)
	if !ok {
		panic(fmt.Sprintf("value of property %q is not a QML list: %#v", property, value))
	}
	return m
}

// Map returns the map value of the named property.
// Map panics if the property is not a map.
func (obj *Common) Map(property string) *Map {
	value := obj.Property(property)
	m, ok := value.(*Map)
	if !ok {
		panic(fmt.Sprintf("value of property %q is not a QML map: %#v", property, value))
	}
	return m
}

// ObjectByName returns the Object value of the descendant object that
// was defined with the objectName property set to the provided value.
// ObjectByName panics if the object is not found.
func (obj *Common) ObjectByName(objectName string) Object {
	obj.assertInitialized()
	cname, cnamelen := util.UnsafeStringData(objectName)
	var dvalue C.DataValue
	var object Object
	RunMain(func() {
		qname := C.newString((*C.char)(cname), C.int(cnamelen))
		defer C.delString(qname)
		C.objectFindChild(obj.addr, qname, &dvalue)
		// unpackDataValue will also initialize the Go type, if necessary.
		value := unpackDataValue(&dvalue, obj.engine)
		if dvalue.dataType == C.DTGoAddr {
			datap := unsafe.Pointer(&dvalue.data)
			fold := (*(**valueFold)(datap))
			if fold.init.IsValid() {
				panic("internal error: custom Go type not initialized")
			}
			cobject := CommonOf(fold.cvalue, fold.engine)
			object = cobject
		} else {
			object, _ = value.(Object)
		}
	})
	if object == nil {
		panic(fmt.Sprintf("cannot find descendant with objectName == %q", objectName))
	}
	return object
}

// Call calls the given object method with the provided parameters.
// Call panics if the method does not exist.
func (obj *Common) Call(method string, params ...interface{}) interface{} {
	obj.assertInitialized()
	if len(params) > len(dataValueArray) {
		panic("too many parameters")
	}
	cmethod, cmethodLen := util.UnsafeStringData(method)
	var result C.DataValue
	var cerr *C.error
	RunMain(func() {
		for i, param := range params {
			packDataValue(param, &dataValueArray[i], obj.engine, jsOwner)
		}
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Panic objectInvoke, %v\n", method)
			}
		}()
		// if obj.addr == nil {
		// fmt.Fprintf(os.Stderr, "objectInvoke, %v %v %v\n", obj.addr, method, len(params))
		// }
		if obj.destroyed {
			// object wasn't destroyed before, checked by assertInitialized, so it was
			// destroyed while waiting to run on the main thread
			// TODO: What to do about this???
		}
		cerr = C.objectInvoke(obj.addr, (*C.char)(cmethod), C.int(cmethodLen), &result, &dataValueArray[0], C.int(len(params)))
	})
	if cerr != nil {
		fmt.Fprintf(os.Stderr, "Common: %#v\n", obj)
	}
	cmust(cerr)
	return unpackDataValue(&result, obj.engine)
}

// Create creates a new instance of the component held by obj.
// The component instance runs under the ctx context. If ctx is nil,
// it runs under the same context as obj.
//
// The Create method panics if called on an object that does not
// represent a QML component.
func (obj *Common) Create(ctx *Context) Object {
	obj.assertInitialized()
	if C.objectIsComponent(obj.addr) == 0 {
		panic("object is not a component")
	}
	var root Common
	root.engine = obj.engine
	RunMain(func() {
		ctxaddr := nilPtr
		if ctx != nil {
			ctxaddr = ctx.addr
		}
		root.setAddr(C.componentCreate(obj.addr, ctxaddr))
	})
	return &root
}

// CreateWindow creates a new instance of the component held by obj,
// and creates a new window holding the instance as its root object.
// The component instance runs under the ctx context. If ctx is nil,
// it runs under the same context as obj.
//
// The CreateWindow method panics if called on an object that
// does not represent a QML component.
func (obj *Common) CreateWindow(ctx *Context) *Window {
	obj.assertInitialized()
	if C.objectIsComponent(obj.addr) == 0 {
		panic("object is not a component")
	}
	var win Window
	win.engine = obj.engine
	RunMain(func() {
		ctxaddr := nilPtr
		if ctx != nil {
			ctxaddr = ctx.addr
		}
		win.setAddr(C.componentCreateWindow(obj.addr, ctxaddr))
	})
	return &win
}

// Destroy finalizes the value and releases any resources used.
// The value must not be used after calling this method.
func (obj *Common) Destroy() {
	// TODO We might hook into the destroyed signal, and prevent this object
	//      from being used in post-destruction crash-prone ways.
	RunMain(func() {
		if obj.addr != nilPtr {
			C.delObjectLater(obj.addr)
			obj.addr = nilPtr
		}
	})
}

var connectedFunction = make(map[*interface{}]bool)

// On connects the named signal from obj with the provided function, so that
// when obj next emits that signal, the function is called with the parameters
// the signal carries.
//
// The provided function must accept a number of parameters that is equal to
// or less than the number of parameters provided by the signal, and the
// resepctive parameter types must match exactly or be conversible according
// to normal Go rules.
//
// For example:
//
//     obj.On("clicked", func() { fmt.Println("obj got a click") })
//
// Note that Go uses the real signal name, rather than the one used when
// defining QML signal handlers ("clicked" rather than "onClicked").
//
// For more details regarding signals and QML see:
//
//     http://qt-project.org/doc/qt-5.0/qtqml/qml-qtquick2-connections.html
//
func (obj *Common) On(signal string, function interface{}) {
	obj.assertInitialized()
	funcv := reflect.ValueOf(function)
	funct := funcv.Type()
	if funcv.Kind() != reflect.Func {
		panic("function provided to On is not a function or method")
	}
	if funct.NumIn() > C.MaxParams {
		panic("function takes too many arguments")
	}
	csignal, csignallen := util.UnsafeStringData(signal)
	var cerr *C.error
	RunMain(func() {
		cerr = C.objectConnect(obj.addr, (*C.char)(csignal), C.int(csignallen), obj.engine.addr, unsafe.Pointer(&function), C.int(funcv.Type().NumIn()))
		if cerr == nil {
			connectedFunction[&function] = true
			stats.connectionsAlive(+1)
		}
	})
	cmust(cerr)
}

//export hookSignalDisconnect
func hookSignalDisconnect(funcp unsafe.Pointer) {
	before := len(connectedFunction)
	delete(connectedFunction, (*interface{})(funcp))
	if before == len(connectedFunction) {
		panic("disconnecting unknown signal function")
	}
	stats.connectionsAlive(-1)
}

//export hookSignalCall
func hookSignalCall(enginep unsafe.Pointer, funcp unsafe.Pointer, args *C.DataValue) {
	engine := engines[enginep]
	if engine == nil {
		panic("signal called after engine was destroyed")
	}
	funcv := reflect.ValueOf(*(*interface{})(funcp))
	funct := funcv.Type()
	numIn := funct.NumIn()
	var params [C.MaxParams]reflect.Value
	for i := 0; i < numIn; i++ {
		arg := (*C.DataValue)(unsafe.Pointer(uintptr(unsafe.Pointer(args)) + uintptr(i)*dataValueSize))
		param := reflect.ValueOf(unpackDataValue(arg, engine))
		if paramt := funct.In(i); param.Type() != paramt {
			// TODO Provide a better error message when this fails.
			param = param.Convert(paramt)
		}
		params[i] = param
	}
	funcv.Call(params[:numIn])
}

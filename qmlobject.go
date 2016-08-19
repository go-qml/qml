package qml

import "image/color"

// Object is the common interface implemented by all QML types.
//
// See the documentation of Common for details about this interface.
type Object interface {
	Common() *Common
	Addr() uintptr
	TypeName() string
	Interface() interface{}
	Set(property string, value interface{})
	Property(name string) interface{}
	Int(property string) int
	Int64(property string) int64
	Float64(property string) float64
	Bool(property string) bool
	String(property string) string
	Color(property string) color.RGBA
	Object(property string) Object
	Map(property string) *Map
	List(property string) *List
	ObjectByName(objectName string) Object
	Call(method string, params ...interface{}) interface{}
	Create(ctx *Context) Object
	CreateWindow(ctx *Context) *Window
	Destroy()
	On(signal string, function interface{})
	Clear()
}

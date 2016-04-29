package qml

// #include <stdlib.h>
//
// #include "capi.h"
//
import "C"

// Context represents a QML context that can hold variables visible
// to logic running within it.
type Context struct {
	Common
}

// SetVar makes the provided value available as a variable with the
// given name for QML code executed within the c context.
//
// If value is a struct, its exported fields are also made accessible to
// QML code as attributes of the named object. The attribute name in the
// object has the same name of the Go field name, except for the first
// letter which is lowercased. This is conventional and enforced by
// the QML implementation.
//
// The engine will hold a reference to the provided value, so it will
// not be garbage collected until the engine is destroyed, even if the
// value is unused or changed.
func (ctx *Context) SetVar(name string, value interface{}) {
	cname, cnamelen := unsafeStringData(name)
	RunMain(func() {
		var dvalue C.DataValue
		packDataValue(value, &dvalue, ctx.engine, cppOwner)

		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)

		C.contextSetProperty(ctx.addr, qname, &dvalue)
	})
}

// SetVars makes the exported fields of the provided value available as
// variables for QML code executed within the c context. The variable names
// will have the same name of the Go field names, except for the first
// letter which is lowercased. This is conventional and enforced by
// the QML implementation.
//
// The engine will hold a reference to the provided value, so it will
// not be garbage collected until the engine is destroyed, even if the
// value is unused or changed.
func (ctx *Context) SetVars(value interface{}) {
	RunMain(func() {
		C.contextSetObject(ctx.addr, wrapGoValue(ctx.engine, value, cppOwner))
	})
}

// Var returns the context variable with the given name.
func (ctx *Context) Var(name string) interface{} {
	cname, cnamelen := unsafeStringData(name)

	var dvalue C.DataValue
	RunMain(func() {
		qname := C.newString(cname, cnamelen)
		defer C.delString(qname)

		C.contextGetProperty(ctx.addr, qname, &dvalue)
	})
	return unpackDataValue(&dvalue, ctx.engine)
}

// Spawn creates a new context that has ctx as a parent.
func (ctx *Context) Spawn() *Context {
	var result Context
	result.engine = ctx.engine
	RunMain(func() {
		result.setAddr(C.contextSpawn(ctx.addr))
	})
	return &result
}

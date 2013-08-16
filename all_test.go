package qml_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/qml"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {}

var _ = Suite(S{})

func (S) SetUpSuite(c *C) {
	qml.Init(nil)
	go qml.Run()
}

type MyStruct struct {
	String  string
	Int     int
	Int64   int64
	Int32   int32
	Float64 float64
	Float32 float32
}

func (S) TestContextSetGetStruct(c *C) {
	//engine := qml.NewEngine()
	//context := engine.RootContext()
	//context.Set("value", &MyStruct{
	//	String:  "<string value>",
	//	Int:     42,
	//	Int64:   42,
	//	Int32:   42,
	//	Float64: 4.2,
	//	Float32: 4.2,
	//})
}

func (S) TestContextSetGetString(c *C) {
	engine := qml.NewEngine()
	context := engine.RootContext()

	context.Set("key", "value")
	c.Assert(context.Get("key"), Equals, "value")

	// XXX Destroy engine.
}

func (S) TestContextSetGetInt64(c *C) {
	engine := qml.NewEngine()
	context := engine.RootContext()

	context.Set("key", int64(42))
	c.Assert(context.Get("key"), Equals, int64(42))

	// XXX Destroy engine.
}

func (S) TestContextSetGetInt32(c *C) {
	engine := qml.NewEngine()
	context := engine.RootContext()

	context.Set("key", int32(42))
	c.Assert(context.Get("key"), Equals, int32(42))

	// XXX Destroy engine.
}

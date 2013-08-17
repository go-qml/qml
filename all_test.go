package qml_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/qml"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	engine *qml.Engine
	context *qml.Context
}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	qml.Init(nil)
	go qml.Run()
}

func (s *S) SetUpTest(c *C) {
	s.engine = qml.NewEngine()
	s.context = s.engine.RootContext()
}

func (s *S) TearDownTest(c *C) {
	s.engine.Close()
}

type MyStruct struct {
	String  string
	Int     int
	Int64   int64
	Int32   int32
	Float64 float64
	Float32 float32
}

func (s *S) TestEngineClosedUse(c *C) {
	s.engine.Close()
	s.engine.Close()
	c.Assert(s.engine.RootContext, PanicMatches, "engine already closed")
}

func (s *S) TestContextSetGetStruct(c *C) {
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

func (s *S) TestContextSetGetString(c *C) {
	s.context.Set("key", "value")
	c.Assert(s.context.Get("key"), Equals, "value")
}

func (s *S) TestContextSetGetInt64(c *C) {
	s.context.Set("key", int64(42))
	c.Assert(s.context.Get("key"), Equals, int64(42))
}

func (s *S) TestContextSetGetInt32(c *C) {
	s.context.Set("key", int32(42))
	c.Assert(s.context.Get("key"), Equals, int32(42))
}

func (s *S) TestContextSetGetFloat64(c *C) {
	s.context.Set("key", float64(42))
	c.Assert(s.context.Get("key"), Equals, float64(42))
}

func (s *S) TestContextSetGetFloat32(c *C) {
	s.context.Set("key", float32(42))
	c.Assert(s.context.Get("key"), Equals, float32(42))
}

func (s *S) TestContextSetGetGoValue(c *C) {
	var value MyStruct
	s.context.Set("key", &value)
	c.Assert(s.context.Get("key"), Equals, &value)
}

func (s *S) TestContextSetObjectGet(c *C) {
	s.context.SetObject(&MyStruct{
		String:  "<string value>",
		Int:     42,
		Int64:   42,
		Int32:   42,
		Float64: 4.2,
		Float32: 4.2,
	})

	c.Assert(s.context.Get("string"), Equals, "<string value>")
	c.Assert(s.context.Get("int64"), Equals, int64(42))
	c.Assert(s.context.Get("int32"), Equals, int32(42))
	c.Assert(s.context.Get("float64"), Equals, float64(4.2))
	c.Assert(s.context.Get("float32"), Equals, float32(4.2))

	v := s.context.Get("int")
	if v != int64(42) && v != int32(42) {
		c.Fatalf("want int32(42) or int64(42), got %T(%v)", v, v)
	}
}

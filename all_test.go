package qml_test

import (
	"fmt"
	. "launchpad.net/gocheck"
	"launchpad.net/qml"
	"runtime"
	"testing"
	"time"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	engine  *qml.Engine
	context *qml.Context
}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	qml.Init(nil)
}

func (s *S) SetUpTest(c *C) {
	qml.SetLogger(c)
	qml.SetStats(true)
	qml.ResetStats()

	stats := qml.GetStats()
	if stats.EnginesAlive > 0 || stats.ValuesAlive > 0 {
		panic(fmt.Sprintf("Test started with values alive: %#v\n", stats))
	}

	s.engine = qml.NewEngine()
	s.context = s.engine.Context()
}

func (s *S) TearDownTest(c *C) {
	s.engine.Destroy()

	retries := 30 // Three seconds top.
	for {
		qml.FlushAll()
		runtime.GC()
		stats := qml.GetStats()
		if stats.EnginesAlive == 0 && stats.ValuesAlive == 0 {
			break
		}
		if retries == 0 {
			panic(fmt.Sprintf("there are objects alive:\n%#v\n", stats))
		}
		retries--
		time.Sleep(100 * time.Millisecond)
		if retries%10 == 0 {
			c.Logf("There are still objects alive; waiting for them to die: %#v\n", stats)
		}
	}

	qml.SetLogger(nil)
}

type MyStruct struct {
	String  string
	True    bool
	False   bool
	Int     int
	Int64   int64
	Int32   int32
	Float64 float64
	Float32 float32
	Any     interface{}
}

var intIs64 bool

func init() {
	var i int = 1<<31 - 1
	intIs64 = (i+1 > 0)
}

func (s *S) TestEngineDestroyedUse(c *C) {
	s.engine.Destroy()
	s.engine.Destroy()
	c.Assert(s.engine.Context, PanicMatches, "engine already destroyed")
}

func (s *S) TestContextGetMissing(c *C) {
	c.Assert(s.context.Get("key"), Equals, nil)
}

func (s *S) TestContextSetGetString(c *C) {
	s.context.Set("key", "value")
	c.Assert(s.context.Get("key"), Equals, "value")
}

func (s *S) TestContextSetGetBool(c *C) {
	s.context.Set("bool", true)
	c.Assert(s.context.Get("bool"), Equals, true)
	s.context.Set("bool", false)
	c.Assert(s.context.Get("bool"), Equals, false)
}

func (s *S) TestContextSetGetInt64(c *C) {
	s.context.Set("key", int64(42))
	c.Assert(s.context.Get("key"), Equals, int64(42))
}

func (s *S) TestContextSetGetInt32(c *C) {
	s.context.Set("key", int32(42))
	c.Assert(s.context.Get("key"), Equals, int32(42))
}

func (s *S) TestContextSetGetInt(c *C) {
	s.context.Set("key", 42)
	if intIs64 {
		c.Assert(s.context.Get("key"), Equals, int64(42))
	} else {
		c.Assert(s.context.Get("key"), Equals, int32(42))
	}
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

func (s *S) TestContextSetGoValueGetProperty(c *C) {
	// This test will touch:
	//
	// - The processing of nesting
	// - Field reading both from a pointer (outter MyStruct) and from a value (inner MyStruct)
	// - Access to an interface{} field (Any)
	// - Proper collection of a JS-owned GoValue wrapper (the result of accessing Any)
	//
	// When changing this test, ensure these tests are covered here or elsewhere.
	value := &MyStruct{Any: MyStruct{String: "<string value>"}}
	s.context.Set("key", &value)

	data := `
		import QtQuick 2.0
		Item{ Component.onCompleted: console.log('string is', key.any.string); }
	`

	component, err := s.engine.Load(qml.String("file.qml", data))
	c.Assert(err, IsNil)

	obj := component.Create(s.context)
	obj.Destroy()

	c.Assert(c.GetTestLog(), Matches, "(?s).*string is <string value>.*")
}

// TODO Test getting of non-existent.

func (s *S) TestContextSetObject(c *C) {
	s.context.SetObject(&MyStruct{
		String:  "<string value>",
		True:    true,
		False:   false,
		Int:     42,
		Int64:   42,
		Int32:   42,
		Float64: 4.2,
		Float32: 4.2,
	})

	c.Assert(s.context.Get("string"), Equals, "<string value>")
	c.Assert(s.context.Get("true"), Equals, true)
	c.Assert(s.context.Get("false"), Equals, false)
	c.Assert(s.context.Get("int64"), Equals, int64(42))
	c.Assert(s.context.Get("int32"), Equals, int32(42))
	c.Assert(s.context.Get("float64"), Equals, float64(4.2))
	c.Assert(s.context.Get("float32"), Equals, float32(4.2))

	if intIs64 {
		c.Assert(s.context.Get("int"), Equals, int64(42))
	} else {
		c.Assert(s.context.Get("int"), Equals, int32(42))
	}
}

func (s *S) TestComponentSetDataError(c *C) {
	_, err := s.engine.Load(qml.String("file.qml", "Item{}"))
	c.Assert(err, ErrorMatches, "file.qml:1 Item is not a type")
}

func (s *S) TestComponentSetData(c *C) {
	const N = 42
	s.context.Set("N", N)
	data := `
		import QtQuick 2.0
		Item { width: N*2; Component.onCompleted: console.log("N is", N) }
	`
	component, err := s.engine.Load(qml.String("file.qml", data))
	c.Assert(err, IsNil)

	pattern := fmt.Sprintf(".* file.qml:3: N is %d\n.*", N)
	c.Assert(c.GetTestLog(), Not(Matches), pattern)

	obj := component.Create(s.context)

	c.Assert(c.GetTestLog(), Matches, pattern)
	c.Assert(obj.Get("width"), Equals, float64(N*2))
}

func (s *S) TestComponentCreateWindow(c *C) {
	data := `
		import QtQuick 2.0
		Item { width: 300; height: 200; }
	`
	component, err := s.engine.Load(qml.String("file.qml", data))
	c.Assert(err, IsNil)

	// TODO How to test this more effectively?
	window := component.CreateWindow(s.context)
	window.Show()
	// Qt doesn't hide the Window if we call it too quickly. :-(
	time.Sleep(100 * time.Millisecond)
	window.Hide()
}

//func (s *S) TestFoo(c *C) {
//	value := MyStruct{String: "<string value>"}
//	s.context.Set("a", &value)
//	s.context.Set("b", &value)
//
//	data := `
//		import QtQuick 2.0
//		Item {
//			Component.onCompleted: {
//				console.log('TEST:', a === b);
//				a = 42;
//			}
//		}
//	`
//
//	component, err := s.engine.Load(qml.String("file.qml", data))
//	c.Assert(err, IsNil)
//	_ = component.Create(s.context)
//
//	c.Assert(s.context.Get("a"), IsNil)
//	c.Assert(c.GetTestLog(), Equals, "")
//}

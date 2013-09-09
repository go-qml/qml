package qml_test

import (
	"fmt"
	. "launchpad.net/gocheck"
	"launchpad.net/qml"
	"runtime"
	"strings"
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
		qml.Flush()
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

type testStruct struct {
	StringValue  string
	TrueValue    bool
	FalseValue   bool
	IntValue     int
	Int64Value   int64
	Int32Value   int32
	Float64Value float64
	Float32Value float32
	AnyValue     interface{}
}

func (ts *testStruct) StringMethod() string {
	return ts.StringValue
}

func (ts *testStruct) PresetString() string {
	ts.StringValue = "<preset value>"
	return ts.StringValue
}

func intIs64() bool {
	var i int = 1<<31 - 1
	return i+1 > 0
}

func intNN(i int) interface{} {
	if intIs64() {
		return int64(i)
	}
	return int32(i)
}

func (s *S) TestEngineDestroyedUse(c *C) {
	s.engine.Destroy()
	s.engine.Destroy()
	c.Assert(s.engine.Context, PanicMatches, "engine already destroyed")
}

var same = "<same>"

var getSetTests = []struct{ set, get interface{} }{
	{"value", same},
	{true, same},
	{false, same},
	{int64(42), same},
	{int32(42), same},
	{float64(42), same},
	{float32(42), same},
	{new(testStruct), same},
	{42, intNN(42)},
}

func (s *S) TestContextGetSet(c *C) {
	for i, t := range getSetTests {
		want := t.get
		if t.get == same {
			want = t.set
		}
		s.context.SetVar("key", t.set)
		c.Assert(s.context.Var("key"), Equals, want,
			Commentf("entry %d is {%v (%T), %v (%T)}", i, t.set, t.set, t.get, t.get))
	}
}

func (s *S) TestContextGetMissing(c *C) {
	c.Assert(s.context.Var("missing"), Equals, nil)
}

func (s *S) TestContextSetGoValueGetProperty(c *C) {
	// This test will touch:
	//
	// - The processing of nesting
	// - Field reading both from a pointer (outter testStruct) and from a value (inner testStruct)
	// - Access to an interface{} field (Any)
	// - Proper collection of a JS-owned GoValue wrapper (the result of accessing Any)
	//
	// When changing this test, ensure these tests are covered here or elsewhere.
	value := &testStruct{AnyValue: testStruct{StringValue: "<string content>"}}
	s.context.SetVar("value", &value)

	data := `
		import QtQuick 2.0
		Item{ Component.onCompleted: console.log('string is', value.anyValue.stringValue); }
	`

	component, err := s.engine.Load(qml.String("file.qml", data))
	c.Assert(err, IsNil)

	inst := component.Create(s.context)
	inst.Destroy()

	c.Assert(c.GetTestLog(), Matches, "(?s).*string is <string content>.*")
}

func (s *S) TestContextSetObject(c *C) {
	s.context.SetVars(&testStruct{
		StringValue:  "<string content>",
		TrueValue:    true,
		FalseValue:   false,
		IntValue:     42,
		Int64Value:   42,
		Int32Value:   42,
		Float64Value: 4.2,
		Float32Value: 4.2,
	})

	c.Assert(s.context.Var("stringValue"), Equals, "<string content>")
	c.Assert(s.context.Var("trueValue"), Equals, true)
	c.Assert(s.context.Var("falseValue"), Equals, false)
	c.Assert(s.context.Var("intValue"), Equals, intNN(42))
	c.Assert(s.context.Var("int64Value"), Equals, int64(42))
	c.Assert(s.context.Var("int32Value"), Equals, int32(42))
	c.Assert(s.context.Var("float64Value"), Equals, float64(4.2))
	c.Assert(s.context.Var("float32Value"), Equals, float32(4.2))
}

func (s *S) TestComponentSetDataError(c *C) {
	_, err := s.engine.Load(qml.String("file.qml", "Item{}"))
	c.Assert(err, ErrorMatches, "file.qml:1 Item is not a type")
}

func (s *S) TestComponentSetData(c *C) {
	const N = 42
	s.context.SetVar("N", N)
	data := `
		import QtQuick 2.0
		Item { width: N*2; Component.onCompleted: console.log("N is", N) }
	`
	component, err := s.engine.Load(qml.String("file.qml", data))
	c.Assert(err, IsNil)

	pattern := fmt.Sprintf(".* file.qml:3: N is %d\n.*", N)
	c.Assert(c.GetTestLog(), Not(Matches), pattern)

	inst := component.Create(s.context)

	c.Assert(c.GetTestLog(), Matches, pattern)
	c.Assert(inst.Field("width"), Equals, float64(N*2))
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

func (s *S) TestObjectIdentity(c *C) {
	value := testStruct{StringValue: "<string content>"}
	s.context.SetVar("a", &value)
	s.context.SetVar("b", &value)

	data := `
		import QtQuick 2.0
		Item {
			Component.onCompleted: {
				console.log('Identical:', a === b);
			}
		}
	`

	component, err := s.engine.Load(qml.String("file.qml", data))
	c.Assert(err, IsNil)
	component.Create(s.context).Destroy()

	c.Assert(c.GetTestLog(), Matches, "(?s).*Identical: true.*")
}

func (s *S) TestRegisterType(c *C) {
	value := &testStruct{StringValue: "new type works!"}
	spec := qml.TypeSpec{
		Location: "GoTest",
		Major:    4,
		Minor:    2,
		Name:     "MyType",
		New:      func() interface{} { return value },
	}
	err := qml.RegisterType(&spec)
	c.Assert(err, IsNil)

	data := `
		import QtQuick 2.0
		import GoTest 4.2
		MyType {
			Component.onCompleted: {
				console.log('Value says:', stringValue)
			}
		}
	`
	component, err := s.engine.Load(qml.String("file.qml", data))
	c.Assert(err, IsNil)

	inst := component.Create(s.context)
	defer inst.Destroy()

	c.Assert(c.GetTestLog(), Matches, "(?s).*Value says: new type works!.*")
}

func (s *S) TestRegisterTypeWriteProperty(c *C) {
	value := &testStruct{}
	spec := qml.TypeSpec{
		Location: "GoTest",
		Major:    4,
		Minor:    2,
		Name:     "NewType",
		New:      func() interface{} { return value },
	}
	qml.RegisterType(&spec)

	data := `
		import GoTest 4.2
		NewType { 
			intValue: 300
			stringValue: "hey"
		}
	`
	component, err := s.engine.Load(qml.String("file.qml", data))
	c.Assert(err, IsNil)

	inst := component.Create(s.context)
	defer inst.Destroy()

	c.Assert(value.StringValue, Equals, "hey")
	c.Assert(value.IntValue, Equals, 300)
}

func (s *S) TestRegisterSingleton(c *C) {
	value := &testStruct{StringValue: "singleton works!"}
	spec := qml.TypeSpec{
		Location: "GoTest",
		Major:    4,
		Minor:    2,
		Name:     "SingletonType",
		New:      func() interface{} { return value },
	}
	err := qml.RegisterSingleton(&spec)
	c.Assert(err, IsNil)

	data := `
		import QtQuick 2.0
		import GoTest 4.2
		Item {
			Component.onCompleted: {
				console.log('Value says:', SingletonType.stringValue)
			}
		}
	`
	component, err := s.engine.Load(qml.String("file.qml", data))
	c.Assert(err, IsNil)

	inst := component.Create(s.context)
	defer inst.Destroy()

	c.Assert(c.GetTestLog(), Matches, "(?s).*Value says: singleton works!.*")
}

func (s *S) TestChanged(c *C) {
	value := &testStruct{StringValue: "<old value>"}
	spec := qml.TypeSpec{
		Location: "GoTest",
		Major:    4,
		Minor:    2,
		Name:     "ChangedType",
		New:      func() interface{} { return value },
	}
	qml.RegisterType(&spec)

	data := `
		import GoTest 4.2
		ChangedType { 
			onStringValueChanged: console.log("String value is now", stringValue)
		}
	`
	component, err := s.engine.Load(qml.String("file.qml", data))
	c.Assert(err, IsNil)

	inst := component.Create(s.context)
	defer inst.Destroy()

	value.StringValue = "<new value>"

	qml.Flush()

	c.Assert(strings.Contains(c.GetTestLog(), "<old value>"), Equals, false)
	c.Assert(strings.Contains(c.GetTestLog(), "<new value>"), Equals, false)

	qml.Changed(value, &value.StringValue)
	qml.Flush()

	c.Assert(strings.Contains(c.GetTestLog(), "<old value>"), Equals, false)
	c.Assert(strings.Contains(c.GetTestLog(), "String value is now <new value>"), Equals, true)
}

func (s *S) TestMethodCall(c *C) {
	value := &testStruct{StringValue: "<string content>"}
	s.context.SetVar("value", value)

	data := `
		import QtQuick 2.0
		Item { Component.onCompleted: console.log('string is', value.stringMethod()); }
	`

	component, err := s.engine.Load(qml.String("file.qml", data))
	c.Assert(err, IsNil)

	inst := component.Create(s.context)
	inst.Destroy()

	c.Assert(c.GetTestLog(), Matches, "(?s).*string is <string content>.*")
}

// TODO presetString is a weird test method, but allows moving forward without
//      methods fully implemented. Change it to something more reasonable once
//      methods work properly.

func (s *S) TestConnectQmlSignalToGoMethod(c *C) {
	value := &testStruct{StringValue: "<string content>"}
	s.context.SetVar("value", value)

	data := `
		import QtQuick 2.0
		Item {
			id: item
			signal testSignal()
			Component.onCompleted: {
				item.testSignal.connect(value.presetString)
				item.testSignal()
			}
		}
	`

	component, err := s.engine.Load(qml.String("file.qml", data))
	c.Assert(err, IsNil)

	inst := component.Create(s.context)
	inst.Destroy()

	c.Assert(value.StringValue, Equals, "<preset value>")
}

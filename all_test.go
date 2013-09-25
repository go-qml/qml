package qml_test

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/niemeyer/qml"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"regexp"
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
	qml.CollectStats(true)
	qml.ResetStats()

	stats := qml.Stats()
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
		// Do not call qml.Flush here. It creates a nested event loop that
		// that attempts to process the deferred object deletes and cannot,
		// because deferred deletes are only processed at the same loop level.
		// So it *reposts* the deferred deletion event, in practice *preventing*
		// these objects from being deleted.
		runtime.GC()
		stats := qml.Stats()
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

type TestType struct {
	private bool // Besides being private, also adds a gap in the reflect field index.

	StringValue  string
	BoolValue    bool
	IntValue     int
	Int64Value   int64
	Int32Value   int32
	Float64Value float64
	Float32Value float32
	AnyValue     interface{}
	ObjectValue  *qml.Object
}

func (ts *TestType) StringMethod() string {
	return ts.StringValue
}

func (ts *TestType) Mod(dividend, divisor int32) (int32, error) {
	if divisor == 0 {
		return 0, fmt.Errorf("<division by zero>")
	}
	return dividend % divisor, nil
}

func (ts *TestType) ChangeString(new string) (old string) {
	old = ts.StringValue
	ts.StringValue = new
	return
}

func (ts *TestType) NotifyStringChanged() {
	qml.Changed(ts, &ts.StringValue)
}

func (ts *TestType) IncrementInt() {
	ts.IntValue++
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
	{new(TestType), same},
	{nil, same},
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

func (s *S) TestContextSetVars(c *C) {
	component, err := s.engine.LoadString("file.qml", "import QtQuick 2.0\nItem { width: 42 }")
	c.Assert(err, IsNil)
	root := component.Create(nil)

	vars := TestType{
		StringValue:  "<content>",
		BoolValue:    true,
		IntValue:     42,
		Int64Value:   42,
		Int32Value:   42,
		Float64Value: 4.2,
		Float32Value: 4.2,
		AnyValue:     nil,
		ObjectValue:  root,
	}
	s.context.SetVars(&vars)

	c.Assert(s.context.Var("stringValue"), Equals, "<content>")
	c.Assert(s.context.Var("boolValue"), Equals, true)
	c.Assert(s.context.Var("intValue"), Equals, intNN(42))
	c.Assert(s.context.Var("int64Value"), Equals, int64(42))
	c.Assert(s.context.Var("int32Value"), Equals, int32(42))
	c.Assert(s.context.Var("float64Value"), Equals, float64(4.2))
	c.Assert(s.context.Var("float32Value"), Equals, float32(4.2))
	c.Assert(s.context.Var("anyValue"), Equals, nil)

	vars.AnyValue = 42
	c.Assert(s.context.Var("anyValue"), Equals, intNN(42))

	c.Assert(s.context.Var("objectValue").(*qml.Object).Int("width"), Equals, 42)
}

func (s *S) TestComponentSetDataError(c *C) {
	_, err := s.engine.LoadString("file.qml", "Item{}")
	c.Assert(err, ErrorMatches, "file:.*/file.qml:1 Item is not a type")
}

func (s *S) TestComponentCreateWindow(c *C) {
	data := `
		import QtQuick 2.0
		Item { width: 300; height: 200; }
	`
	component, err := s.engine.LoadString("file.qml", data)
	c.Assert(err, IsNil)

	// TODO How to test this more effectively?
	window := component.CreateWindow(nil)
	window.Show()
	// Qt doesn't hide the Window if we call it too quickly. :-(
	time.Sleep(100 * time.Millisecond)
	window.Hide()
}

type TestData struct {
	*C
	engine    *qml.Engine
	context   *qml.Context
	component *qml.Object
	root      *qml.Object
	value     *TestType
}

var tests = []struct {
	Summary string
	Value   TestType

	Init func(d *TestData)

	// The QML provided is run with the initial state above, and
	// then checks are made to ensure the provided state is found.
	QML      string
	QMLLog   string
	QMLValue TestType

	// The function provided is run with the post-QML state above,
	// and then checks are made to ensure the provided state is found.
	Done      func(d *TestData)
	DoneLog   string
	DoneValue TestType
}{
	{
		Summary: "Setting and reading of context variables",
		Value:   TestType{StringValue: "<content>", IntValue: 42},
		QML: `
			Item {
				Component.onCompleted: {
					console.log("String is", value.stringValue)
					console.log("Int is", value.intValue)
					console.log("Any is", value.anyValue)
				}
			}
		`,
		QMLLog: "String is <content>.*Int is 42.*Any is undefined",
	},
	{
		Summary: "Reading of nested field via a value (not pointer) in an interface",
		Value:   TestType{AnyValue: TestType{StringValue: "<content>"}},
		QML:     `Item { Component.onCompleted: console.log("String is", value.anyValue.stringValue) }`,
		QMLLog:  "String is <content>",
	},
	{
		Summary: "Reading of value fields",
		QML: `
			Item {
				width: 123;
				property string s: "foo";
			}
		`,
		Done: func(d *TestData) {
			d.Check(d.root.Int("width"), Equals, 123)
			d.Check(d.root.String("s"), Equals, "foo")
		},
	},
	{
		Summary: "Read object properties",
		QML: `
			Item {
				property bool boolp: true
				property int intp: 1
				property var int64p: 4294967296
				property real float32p: 1.1
				property double float64p: 1.1
				property string stringp: "<content>"
				property var objectp: Rectangle { width: 123 }
				property var nilp: null
			}
		`,
		Done: func(d *TestData) {
			obj := d.root
			d.Check(obj.Bool("boolp"), Equals, true)
			d.Check(obj.Int("intp"), Equals, 1)
			d.Check(obj.Int64("intp"), Equals, int64(1))
			d.Check(obj.Int64("int64p"), Equals, int64(4294967296))
			d.Check(obj.Float64("intp"), Equals, float64(1))
			d.Check(obj.Float64("int64p"), Equals, float64(4294967296))
			d.Check(obj.Float64("float32p"), Equals, float64(1.1))
			d.Check(obj.Float64("float64p"), Equals, float64(1.1))
			d.Check(obj.String("stringp"), Equals, "<content>")
			d.Check(obj.Object("objectp").Int("width"), Equals, 123)
			d.Check(obj.Property("nilp"), Equals, nil)

			d.Check(func() { obj.Bool("intp") }, Panics, `value of property "intp" is not a bool: 1`)
			d.Check(func() { obj.Int("boolp") }, Panics, `value of property "boolp" cannot be represented as an int: true`)
			d.Check(func() { obj.Int64("boolp") }, Panics, `value of property "boolp" cannot be represented as an int64: true`)
			d.Check(func() { obj.Float64("boolp") }, Panics, `value of property "boolp" cannot be represented as a float64: true`)
			d.Check(func() { obj.String("boolp") }, Panics, `value of property "boolp" is not a string: true`)
			d.Check(func() { obj.Object("boolp") }, Panics, `value of property "boolp" is not a *qml.Object: true`)
			d.Check(func() { obj.Property("missing") }, Panics, `object does not have a "missing" property`)
		},
	},
	{
		Summary: "No access to private fields",
		Value:   TestType{private: true},
		QML:     `Item { Component.onCompleted: console.log("Private is", value.private); }`,
		QMLLog:  "Private is undefined",
	},
	{
		Summary: "Setting of value fields",
		QML: `
			Item {
				property var obj: null

				onObjChanged:     console.log("String is", obj.stringValue)
				onWidthChanged:   console.log("Width is", width)
				onHeightChanged:  console.log("Height is", height)
			}
		`,
		Done: func(d *TestData) {
			value := TestType{StringValue: "<content>"}
			d.root.Set("obj", &value)
			d.root.Set("width", 300)
			d.root.Set("height", 200)
		},
		DoneLog: "String is <content>.*Width is 300.*Height is 200",
	},
	{
		Summary: "Identical values remain identical when possible",
		Init: func(d *TestData) {
			d.context.SetVar("a", d.value)
			d.context.SetVar("b", d.value)
		},
		QML:    `Item { Component.onCompleted: console.log('Identical:', a === b); }`,
		QMLLog: "Identical: true",
	},
	{
		Summary: "Object finding via objectName",
		QML:     `Item { Item { objectName: "subitem"; property string s: "<found>" } }`,
		Done: func(d *TestData) {
			obj := d.root.ObjectByName("subitem")
			d.Check(obj.String("s"), Equals, "<found>")
			d.Check(func() { d.root.ObjectByName("foo") }, Panics, `cannot find descendant with objectName == "foo"`)
		},
	},
	{
		Summary: "Register Go type",
		Value:   TestType{StringValue: "<content>"},
		QML: `
			import GoTest 4.2
			GoType { Component.onCompleted: console.log("String is", stringValue) }
		`,
		QMLLog: "String is <content>",
	},
	{
		Summary: "Write Go type property",
		QML: `
			import GoTest 4.2
			GoType { stringValue: "<new>"; intValue: 300 }
		`,
		QMLValue: TestType{StringValue: "<new>", IntValue: 300},
	},
	{
		Summary: "Singleton type registration",
		Value:   TestType{StringValue: "<content>"},
		QML: `
			import GoTest 4.2
			Item { Component.onCompleted: console.log("String is", GoSingleton.stringValue) }
		`,
		QMLLog: "String is <content>",
	},
	{
		Summary: "qml.Changed triggers a QML slot",
		Value:   TestType{StringValue: "<old>"},

		QML: `
			import GoTest 4.2
			GoType { onStringValueChanged: console.log("String is", stringValue) }
		`,
		QMLLog:   "!String is",
		QMLValue: TestType{StringValue: "<old>"},

		Done: func(d *TestData) {
			d.value.StringValue = "<new>"
			qml.Changed(d.value, &d.value.StringValue)
		},
		DoneLog:   "String is <new>",
		DoneValue: TestType{StringValue: "<new>"},
	},
	{
		Summary: "qml.Changed must not trigger on the wrong field",
		Value:   TestType{StringValue: "<old>"},
		QML: `
			import GoTest 4.2
			GoType { onStringValueChanged: console.log("String is", stringValue) }
		`,
		Done: func(d *TestData) {
			d.value.StringValue = "<new>"
			qml.Changed(d.value, &d.value.IntValue)
		},
		DoneLog: "!String is",
	},
	{
		Summary: "qml.Changed triggers multiple wrappers of the same value",
		Value:   TestType{StringValue: "<old>"},
		Init: func(d *TestData) {
			d.context.SetVar("v1", d.value)
			d.context.SetVar("v2", d.value)
			d.context.SetVar("v3", d.value)
		},

		QML: `
			import GoTest 4.2
			Item {
				property var p1: GoType { onStringValueChanged: console.log("p1 has", stringValue) }
				property var p2: GoType { onStringValueChanged: console.log("p2 has", stringValue) }
				property var p3: GoType { onStringValueChanged: console.log("p3 has", stringValue) }
				Connections { target: v1; onStringValueChanged: console.log("v1 has", v1.stringValue) }
				Connections { target: v2; onStringValueChanged: console.log("v2 has", v2.stringValue) }
				Connections { target: v3; onStringValueChanged: console.log("v3 has", v3.stringValue) }
			}
		`,
		QMLLog:   "![pv][123] has <old>",
		QMLValue: TestType{StringValue: "<old>"},

		Done: func(d *TestData) {
			d.value.StringValue = "<new>"
			qml.Changed(d.value, &d.value.StringValue)
		},
		// Why are v3-v1 reversed? Is QML registering connections in reversed order?
		DoneLog: "v3 has <new>.*v2 has <new>.*v1 has <new>.*p1 has <new>.*p2 has <new>.*p3 has <new>.*",
	},
	{
		Summary: "qml.Changed updates bindings",
		Value:   TestType{StringValue: "<old>"},
		QML:     `Item { property string s: "String is " + value.stringValue }`,
		Done: func(d *TestData) {
			d.value.StringValue = "<new>"
			qml.Changed(d.value, &d.value.StringValue)
			d.Check(d.root.String("s"), Equals, "String is <new>")
		},
	},
	{
		Summary:  "Call a Go method without arguments or result",
		Value:    TestType{IntValue: 42},
		QML:      `Item { Component.onCompleted: console.log("Undefined is", value.incrementInt()); }`,
		QMLLog:   "Undefined is undefined",
		QMLValue: TestType{IntValue: 43},
	},
	{
		Summary:  "Call a Go method with one argument and one result",
		Value:    TestType{StringValue: "<old>"},
		QML:      `Item { Component.onCompleted: console.log("String was", value.changeString("<new>")); }`,
		QMLLog:   "String was <old>",
		QMLValue: TestType{StringValue: "<new>"},
	},
	{
		Summary: "Call a Go method with multiple results",
		QML: `
			Item {
				Component.onCompleted: {
					var r = value.mod(42, 4);
					console.log("mod is", r[0], "and err is", r[1]);
				}
			}
		`,
		QMLLog: `mod is 2 and err is undefined`,
	},
	{
		Summary: "Call a Go method that returns an error",
		QML: `
			Item {
				Component.onCompleted: {
					var r = value.mod(0, 0);
					console.log("err is", r[1].error());
				}
			}
		`,
		QMLLog: `err is <division by zero>`,
	},
	{
		Summary: "Call a Go method that recurses back into the GUI thread",
		QML: `
			Item {
				Connections {
					target: value
					onStringValueChanged: console.log("Notification arrived")
				}
				Component.onCompleted: {
					value.notifyStringChanged()
				}
			}
		`,
		QMLLog: "Notification arrived",
	},
	{
		Summary: "Connect a QML signal to a Go method",
		Value:   TestType{StringValue: "<old>"},
		QML: `
			Item {
				id: item
				signal testSignal(string s)
				Component.onCompleted: {
					item.testSignal.connect(value.changeString)
					item.testSignal("<new>")
				}
			}
		`,
		QMLValue: TestType{StringValue: "<new>"},
	},
	{
		Summary: "Call a QML method with no result or parameters from Go",
		QML:     `Item { function f() { console.log("f was called"); } }`,
		Done:    func(d *TestData) { d.Check(d.root.Call("f"), IsNil) },
		DoneLog: "f was called",
	},
	{
		Summary: "Call a QML method with result and parameters from Go",
		QML:     `Item { function add(a, b) { return a+b; } }`,
		Done:    func(d *TestData) { d.Check(d.root.Call("add", 1, 2), Equals, int32(3)) },
	},
	{
		Summary: "Call a QML method with a custom type",
		Value:   TestType{StringValue: "<content>"},
		QML:     `Item { function log(value) { console.log("String is", value.stringValue) } }`,
		Done:    func(d *TestData) { d.root.Call("log", d.value) },
		DoneLog: "String is <content>",
	},
	{
		Summary: "Call a QML method that returns a QML object",
		QML: `
			Item {
				property var custom: Rectangle { width: 300; }
				function f() { return custom }
			}
		`,
		Done: func(d *TestData) {
			d.Check(d.root.Call("f").(*qml.Object).Int("width"), Equals, 300)
		},
	},
	{
		Summary: "Call a QML method that holds a custom type past the return point",
		QML: `
			Item {
				property var held
				function hold(v) { held = v; gc(); gc(); }
				function log()   { console.log("String is", held.stringValue) }
			}`,
		Done: func(d *TestData) {
			value := TestType{StringValue: "<content>"}
			stats := qml.Stats()
			d.root.Call("hold", &value)
			d.Check(qml.Stats().ValuesAlive, Equals, stats.ValuesAlive+1)
			d.root.Call("log")
			d.root.Call("hold", nil)
			d.Check(qml.Stats().ValuesAlive, Equals, stats.ValuesAlive)
		},
		DoneLog: "String is <content>",
	},
	{
		Summary: "Ensure URL of provided file is correct by loading a local file",
		Init: func(d *TestData) {
			data, err := base64.StdEncoding.DecodeString("R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==")
			d.Assert(err, IsNil)
			err = ioutil.WriteFile("test.gif", data, 0644)
			d.Check(err, IsNil)
		},
		QML:    `Image { source: "test.gif"; Component.onCompleted: console.log("Ready:", status == Image.Ready) }`,
		QMLLog: "Ready: true",
		Done:   func(d *TestData) { os.Remove("test.gif") },
	},
	{
		Summary: "Create window with non-window root object",
		QML:     `Rectangle { width: 300; height: 200; function inc(x) { return x+1 } }`,
		Done: func(d *TestData) {
			win := d.component.CreateWindow(nil)
			root := win.Root()
			d.Check(root.Int("width"), Equals, 300)
			d.Check(root.Int("height"), Equals, 200)
			d.Check(root.Call("inc", 42), Equals, int32(43))
			root.Destroy()
		},
	},
	{
		Summary: "Create window with window root object",
		QML:     `
			import QtQuick.Window 2.0
			Window { title: "<title>"; width: 300; height: 200 }
		`,
		Done: func(d *TestData) {
			win := d.component.CreateWindow(nil)
			root := win.Root()
			d.Check(root.String("title"), Equals, "<title>")
			d.Check(root.Int("width"), Equals, 300)
			d.Check(root.Int("height"), Equals, 200)
		},
	},
	{
		Summary: "Pass a *Value back into a method",
		QML:     `Rectangle { width: 300; function log(r) { console.log("Width is", r.width) } }`,
		Done:    func(d *TestData) { d.root.Call("log", d.root) },
		DoneLog: "Width is 300",
	},
	{
		Summary: "Create a QML-defined component in Go",
		QML:     `Item { property var comp: Component { Rectangle { width: 300 } } }`,
		Done: func(d *TestData) {
			rect := d.root.Object("comp").Create(nil)
			d.Check(rect.Int("width"), Equals, 300)
			d.Check(func() { d.root.Create(nil) }, Panics, "object is not a component")
			d.Check(func() { d.root.CreateWindow(nil) }, Panics, "object is not a component")
		},
	},
	{
		Summary: "Call a Qt method that has no result",
		QML:     `Item { Component.onDestruction: console.log("item destroyed") }`,
		Done: func(d *TestData) {
			// Create a local instance to avoid double-destroying it.
			root := d.component.Create(nil)
			root.Call("deleteLater")
			time.Sleep(100 * time.Millisecond)
		},
		DoneLog: "item destroyed",
	},
	{
		Summary: "Connect to a QML signal",
		QML: `
			Item {
				id: item
				signal doIt()
				function emitDoIt() { item.doIt() }
			}
		`,
		Done: func(d *TestData) {
			itWorks := false
			d.root.On("doIt", func() { itWorks = true })
			d.Check(itWorks, Equals, false)
			d.root.Call("emitDoIt")
			d.Check(itWorks, Equals, true)
			d.Check(func() { d.root.On("missing", func() {}) }, Panics, `object has no "missing" signal`)
		},
	},
}

var tablef = flag.String("tablef", "", "if provided, TestTable only runs tests with a summary matching the regexp")

func (s *S) TestTable(c *C) {
	var goTypeValue *TestType = &TestType{}

	typeSpec := qml.TypeSpec{
		Location: "GoTest",
		Major:    4,
		Minor:    2,
		Name:     "GoType",
		New:      func() interface{} { return goTypeValue },
	}
	err := qml.RegisterType(&typeSpec)
	c.Assert(err, IsNil)

	singletonSpec := qml.TypeSpec{
		Location: "GoTest",
		Major:    4,
		Minor:    2,
		Name:     "GoSingleton",
		New:      func() interface{} { return goTypeValue },
	}
	err = qml.RegisterSingleton(&singletonSpec)
	c.Assert(err, IsNil)

	filter := regexp.MustCompile("")
	if tablef != nil {
		filter = regexp.MustCompile(*tablef)
	}

	for i := range tests {
		s.TearDownTest(c)
		t := &tests[i]
		header := fmt.Sprintf("----- Running table test %d: %s -----", i, t.Summary)
		if !filter.MatchString(header) {
			continue
		}
		c.Log(header)
		s.SetUpTest(c)

		value := t.Value
		goTypeValue = &value
		s.context.SetVar("value", &value)

		testData := TestData{
			C:       c,
			value:   &value,
			engine:  s.engine,
			context: s.context,
		}

		if t.Init != nil {
			t.Init(&testData)
			if c.Failed() {
				c.FailNow()
			}
		}

		component, err := s.engine.LoadString("file.qml", "import QtQuick 2.0\n"+strings.TrimSpace(t.QML))
		c.Assert(err, IsNil)

		logMark := c.GetTestLog()

		// The component instance is destroyed before the loop ends below,
		// but do a defer to ensure it will be destroyed if the test fails.
		root := component.Create(nil)
		defer root.Destroy()

		testData.component = component
		testData.root = root

		if t.QMLLog != "" {
			logged := c.GetTestLog()[len(logMark):]
			if t.QMLLog[0] == '!' {
				c.Check(logged, Not(Matches), "(?s).*"+t.QMLLog[1:]+".*")
			} else {
				c.Check(logged, Matches, "(?s).*"+t.QMLLog+".*")
			}
		}

		if t.QMLValue != (TestType{}) {
			c.Check(value.StringValue, Equals, t.QMLValue.StringValue)
			c.Check(value.BoolValue, Equals, t.QMLValue.BoolValue)
			c.Check(value.IntValue, Equals, t.QMLValue.IntValue)
			c.Check(value.Int64Value, Equals, t.QMLValue.Int64Value)
			c.Check(value.Int32Value, Equals, t.QMLValue.Int32Value)
			c.Check(value.Float64Value, Equals, t.QMLValue.Float64Value)
			c.Check(value.Float32Value, Equals, t.QMLValue.Float32Value)
			c.Check(value.AnyValue, Equals, t.QMLValue.AnyValue)
		}

		if !c.Failed() {
			logMark := c.GetTestLog()

			if t.Done != nil {
				t.Done(&testData)
			}

			if t.DoneLog != "" {
				logged := c.GetTestLog()[len(logMark):]
				if t.DoneLog[0] == '!' {
					c.Check(logged, Not(Matches), "(?s).*"+t.DoneLog[1:]+".*")
				} else {
					c.Check(logged, Matches, "(?s).*"+t.DoneLog+".*")
				}
			}

			if t.DoneValue != (TestType{}) {
				c.Check(value.StringValue, Equals, t.DoneValue.StringValue)
				c.Check(value.BoolValue, Equals, t.DoneValue.BoolValue)
				c.Check(value.IntValue, Equals, t.DoneValue.IntValue)
				c.Check(value.Int64Value, Equals, t.DoneValue.Int64Value)
				c.Check(value.Int32Value, Equals, t.DoneValue.Int32Value)
				c.Check(value.Float64Value, Equals, t.DoneValue.Float64Value)
				c.Check(value.Float32Value, Equals, t.DoneValue.Float32Value)
				c.Check(value.AnyValue, Equals, t.DoneValue.AnyValue)
			}
		}

		root.Destroy()

		if c.Failed() {
			c.FailNow() // So relevant logs are at the bottom.
		}
	}
}

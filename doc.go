// Package qml offers graphical QML application support for the Go language.
//
// Warning
//
// This package is in an alpha stage, and still in heavy development. APIs may
// change, and things may break.
//
// At this time contributors and developers that are interested in tracking the
// development closely are encouraged to use it. If you'd prefer a more stable
// release, please hold on a bit and subscribe to the mailing list for news. It's
// in a pretty good state, so it shall not take too long.
//
// See http://github.com/niemeyer/qml for details.
//
// Introduction
//
// The qml package enables Go programs to display and manipulate graphical content
// using Qt's QML framework. QML uses a declarative language to express structure
// and style, and supports JavaScript for in-place manipulation of the described
// content. When using the Go qml package, such QML content can also interact with
// Go values, making use of its exported fields and methods, and even explicitly
// creating new instances of registered Go types.
//
// A simple Go application that integrates with QML may perform the following steps
// for offering a graphical interface:
//
//   * Initialize the qml package (see Init)
//   * Create an engine for loading and running QML content (see NewEngine)
//   * Make Go values and types available to QML (see Context.SetVar and RegisterType)
//   * Load QML content (see Engine.LoadString and Engine.LoadFile)
//   * Create a new window for the content (see Component.CreateWindow)
//   * Show the window and wait for it to be closed (see Window.Show and Window.Wait)
//
// Some of these topics are covered below, and may also be observed in practice
// in the following examples:
//
//   https://github.com/niemeyer/qml/tree/master/examples
//
// Making Go values available to QML
//
// The simplest way of making a Go value available to QML code is setting it
// as a variable of the engine's root context, as in:
//
//     context := engine.Context()
//     context.SetVar("person", &Person{Name: "Ale"})
//
// This logic would enable the following QML code to successfully run:
//
//     import QtQuick 2.0
//     Item {
//         Component.onCompleted: console.log("Name is", person.name)
//     }
//
// While this method is a quick way to get started, it is also fairly limited.
// For more flexibility, a Go type may be registered so that QML code can
// natively create new instances in an arbitrary position of the structure.
// This may be achieved by simply registering the desired type via the RegisterType
// function:
//
//    qml.RegisterTypes("GoExtensions", 1, 0, []qml.TypeSpec{{
//        Name: "Person",
//        Init: func(p *Person, obj qml.Object) { p.Name = "<none>" },
//    }})
//
// With this in place, QML code can create new instances of Person by itself:     
//
//    import QtQuick 2.0
//    import GoExtensions 1.0
//    Item{
//        Person {
//            id: person
//            name: "Ale"
//        }
//        Component.onCompleted: console.log("Name is", person.name)
//    }
// 
// Using either mechanism, the methods and fields from Go values are available
// to QML logic as methods and properties of the respective QML object. As
// required by QML, the Go method and field names are lowercased according to
// the following scheme when being accesed from QML:
//
//     value.Name      => value.name
//     value.UPPERName => value.upperName
//     value.UPPER     => value.upper
//
// Setters and getters
//
// In addition to directly reading and writing value fields from QML code, as
// described above, Go values may also intercept writes to specific fields by
// declaring a setter method according to common Go conventions.
//
// For example:
//
//     type Person struct {
//         Name string
//     }
//
//     func (p *Person) SetName(name string) {
//         fmt.Println("Old name is", p.Name)
//         p.Name = name
//         fmt.Println("New name is", p.Name)
//     }
//
// In the example above, whenever QML logic writes to the Person.Name field via
// any means (including object declarations) the SetName method is invoked.
//
// A setter method may also be used in conjunction with a getter method rather
// than a real type field. A method is only considered a getter in the presence
// of the respective setter, and according to common Go conventions it must not
// have the Get prefix.
//
// Inside QML logic, the getter and setter pair is seen as a single object property.
//
// For example:
//
//     type Person struct{}
//
//     func (p *Person) Name() string {
//         return p.loadName()
//     }
//
//     func (p *Person) SetName(name string) {
//         p.saveName(name)
//     }
//
// The type above could be used within QML as follows:
//
//     import GoExtensions 1.0
//     Person {
//         id: person
//         name: "Ale"
//         Component.onCompleted: console.log("Name is", person.name)
//     }
//
// Painting
//
// Custom types implemented in Go may have displayable content by defining
// a Paint method such as:
//
//     func (p *Person) Paint(painter *qml.Painter) {
//         // ... OpenGL calls with the github.com/niemeyer/qml/gl package ...
//     }
//
// A simple example is available at:
//
//     https://github.com/niemeyer/qml/tree/master/examples/painting
//
package qml

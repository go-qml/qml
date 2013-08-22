import QtQuick 2.0
//import GoTypes 1.0

Rectangle {
  width: 300; height: 150
  Text {
    id: text
    x: 0
    y: 0
    text: "Go says: " + message.text
    //MyStruct {
    //	//foo: "text"
    //}
    SequentialAnimation {
      running: true
      NumberAnimation { target: text; property: "x"; to: 70; duration: 500 }
      NumberAnimation { target: text; property: "y"; to: 70; duration: 500 }
      ScriptAction {
        //script: { console.log("Setting obj.bar"); obj.bar = "<new value for bar>"; text.blah = "caca"; text.slotName(); }
      }
      PauseAnimation { duration: 500 }
      ScriptAction {
        script: { console.log("After pause, Go says:", message.text) }
        //script: { console.warn("Setting obj.foo"); obj.foo = "<new value for foo>"; }
      }
    }
    //property string blah
    //onBlahChanged: console.log("blah is now " + blah)
    //Connections {
    //  target: obj
    //  onBarChanged: console.log("obj.bar changed!")
    //  onFooChanged: console.log("obj.foo changed!")
    //}
    //function slotName() {
    //    console.log("slotName called!")
    //}
  }
}

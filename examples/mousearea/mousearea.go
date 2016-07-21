package main

import (
	"fmt"
	"os"

	"gopkg.in/qml.v1"
)

const test_qml = `
import QtQuick 2.2
import QtQuick.Controls 1.1

ApplicationWindow {
    width: 640
    height: 280

    MouseArea {
        objectName: "mouseArea"
        hoverEnabled: true

        anchors.fill: parent
    }
}
`

func main() {
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	engine := qml.NewEngine()
	component, err := engine.LoadString("test.qml", test_qml)
	if err != nil {
		return err
	}

	win := component.CreateWindow(nil)

	mouseArea := win.Root().ObjectByName("mouseArea")
	mouseArea.On("positionChanged", func(mouseEvent qml.Object) {
		fmt.Printf("X %d Y %d\n", mouseEvent.Int("x"), mouseEvent.Int("y"))
	})

	win.Show()
	win.Wait()

	return nil
}

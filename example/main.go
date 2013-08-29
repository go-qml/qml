package main

import (
	"launchpad.net/qml"
)

type Message struct {
	Text string
}

func main() {
	qml.Init(nil)
	engine := qml.NewEngine()
	component, err := engine.Load(qml.File("example.qml"))
	if err != nil {
		panic(err)
	}

	context := engine.Context()
	context.Set("message", &Message{"Hello from Go!"})

	window := component.CreateWindow(nil)
	window.Show()
	window.Wait()
}

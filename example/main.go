package main

import (
	"io/ioutil"
	"launchpad.net/qml"
)

type Message struct {
	Text string
}

func main() {
	data, err := ioutil.ReadFile("example.qml")
	if err != nil {
		panic(err)
	}

	qml.Init(nil)
	engine := qml.NewEngine()
	component := qml.NewComponent(engine)
	err = component.SetData("example.qml", data)
	if err != nil {
		panic(err)
	}

	context := engine.Context()
	context.Set("message", &Message{"Hello from Go!"})

	window := component.CreateWindow(engine.RootContext())
	window.Show()

	select{}
}

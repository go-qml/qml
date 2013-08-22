package main

import (
	"io/ioutil"
	"launchpad.net/qml"
	"time"
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

	context := engine.RootContext()
	context.Set("message", &Message{"Hello from Go!"})

	window := component.CreateWindow(engine.RootContext())
	window.Show()

	for i := 0; i < 10; i++ {
		go time.Sleep(5)
	}

	time.Sleep(10)

	qml.Run()
}

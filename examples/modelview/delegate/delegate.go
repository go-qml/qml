package main

import (
	"fmt"
	"github.com/niemeyer/qml"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	qml.Init(nil)
	engine := qml.NewEngine()
	colors := &Colors{[]string{"red", "green", "blue", "black"}}
	engine.Context().SetVar("colors", colors)
	component, err := engine.LoadFile("delegate.qml")
	if err != nil {
		return err
	}
	window := component.CreateWindow(nil)
	window.Show()
	window.Wait()
	return nil
}

type Colors struct {
	names []string
}

func (colors *Colors) Len() int {
	return len(colors.names)
}

func (colors *Colors) Name(index int32) string {
	return colors.names[index]
}

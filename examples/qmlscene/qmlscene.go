package main

import (
	"fmt"
	"os"

	"github.com/limetext/qml-go"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <qml file>\n", os.Args[0])
		os.Exit(1)
	}
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	engine := qml.NewEngine()

	engine.On("quit", func() { os.Exit(0) })

	component, err := engine.LoadFile(os.Args[1])
	if err != nil {
		return err
	}
	window := component.CreateWindow(nil)
	window.Show()
	window.Wait()
	return nil
}

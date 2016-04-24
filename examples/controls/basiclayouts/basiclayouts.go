package main

import (
	"fmt"
	"os"

	"github.com/limetext/qml-go"
)

func main() {
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	engine := qml.NewEngine()

	controls, err := engine.LoadFile("main.qml")
	if err != nil {
		return err
	}

	window := controls.CreateWindow(nil)

	window.Show()
	window.Wait()
	return nil
}

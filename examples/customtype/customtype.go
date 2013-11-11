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

type GoType struct {
	Text string
}

func (v *GoType) OnTextChanged() {
	fmt.Println("Text changed...")
}

func run() error {
	qml.Init(nil)

	qml.RegisterTypes("GoExtensions", 1, 0, []qml.TypeSpec{{
		Name: "GoType",
		New: func() interface{} { return &GoType{} },
	}})

	engine := qml.NewEngine()
	component, err := engine.LoadFile("customtype.qml")
	if err != nil {
		return err
	}

	value := component.Create(nil)
	fmt.Println("Text is:", value.Interface().(*GoType).Text)

	return nil
}

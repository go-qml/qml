package main

import (
	"fmt"

	"gopkg.in/qml.v1"
)

func main() {
	err := qml.Run(run)
	if err != nil {
		panic(err)
	}
}

type GoStruct struct {
}

func (gs *GoStruct) ReturnGoType() *GoStruct {
	return gs
}
func (gs *GoStruct) UseGoType(v *GoStruct) {
	fmt.Println("Successfully called UseGoType()")
}

func run() error {
	engine := qml.NewEngine()
	context := engine.Context()
	context.SetVar("gostruct", &GoStruct{})
	component, err := engine.LoadFile("main.qml")
	if err != nil {
		return err
	}
	win := component.CreateWindow(nil)
	win.Show()
	win.Wait()
	return nil
}

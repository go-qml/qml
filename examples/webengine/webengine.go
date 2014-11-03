package main

import (
	"fmt"
	"os"

	"gopkg.in/qml.v1"
	"gopkg.in/qml.v1/webengine"
)

func main() {
	fmt.Println(qml.Run(func() error {
		webengine.Initialize()

		engine := qml.NewEngine()
		engine.On("quit", func() { os.Exit(0) })

		component, err := engine.LoadFile("webengine.qml")
		if err != nil {
			return err
		}
		win := component.CreateWindow(nil)
		win.Show()
		win.Wait()
		return nil
	}))
}

package main

import (
	"fmt"
	"github.com/niemeyer/qml"
	"image/color"
	"math/rand"
	"os"
	"time"
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
	colors := &Colors{}
	engine.Context().SetVar("colors", colors)
	component, err := engine.LoadFile("delegate.qml")
	if err != nil {
		return err
	}
	window := component.CreateWindow(nil)
	window.Show()
	n := func() uint8 { return uint8(rand.Intn(256)) }
	for i := 0; i < 5; i++ {
		colors.Add(color.RGBA{n(), n(), n(), 0xff})
	}
	go func() {
		for i := 0; i < 10; i++ {
			for i := 0; i < 5; i++ {
				colors.Change(i, color.RGBA{n(), n(), n(), 0xff})
				time.Sleep(1 * time.Second)
			}
		}
	}()
	window.Wait()
	return nil
}

type Colors struct {
	list []*ColorItem
	Len  int
}

type ColorItem struct {
	Color color.RGBA
}

func (colors *Colors) Add(c color.RGBA) {
	colors.list = append(colors.list, &ColorItem{c})
	colors.Len = len(colors.list)
	qml.Changed(colors, &colors.Len)
}

func (colors *Colors) Color(index int) *ColorItem {
	return colors.list[index]
}

func (colors *Colors) Change(index int, c color.RGBA) {
	item := colors.list[index]
	item.Color = c
	qml.Changed(item, &item.Color)
}

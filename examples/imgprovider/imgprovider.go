package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"bitbucket.org/kardianos/osext"
	"gopkg.in/v0/qml"
)

// fatal puts formatted msg on std error and exits.
func fatal(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	os.Exit(1)
}

func main() {
	dir, e := osext.ExecutableFolder()
	if nil != e {
		fatal("osext cannot find executable folder: %v\n", e)
	}
	// Changing directory ensures the program can be run from any directory.
	// Otherwise it cannot find the qml file.

	e = os.Chdir(dir)
	if nil != e {
		fatal("cannot change wd to '%s': %v\n", dir, e)
	}

	if err := run(); err != nil {
		fatal("error: %v\n", err)
	}
}

func run() error {
	qml.Init(nil)
	engine := qml.NewEngine()
	engine.AddImageProvider("pwd", func(id string, width, height int) image.Image {
		f, err := os.Open(id)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		image, err := png.Decode(f)
		if err != nil {
			panic(err)
		}
		return image
	})

	component, err := engine.LoadFile("imgprovider.qml")
	if err != nil {
		return err
	}

	win := component.CreateWindow(nil)
	win.Show()
	win.Wait()

	return nil
}

package main

import (
	"fmt"
	"gopkg.in/v0/qml"
	"image/color"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type GoType struct {
	My_go_text string
	//add some other valid types here
	//as described here:
	//http://godoc.org/github.com/niemeyer/qml#Object
        //Int(property string) int
	My_go_int int
        //Int64(property string) int64
	My_go_int64 int64
        //Float64(property string) float64
	My_go_float64 float64
        //Bool(property string) bool
	My_go_bool bool
        //String(property string) string
	My_go_string string
        //Color(property string) color.RGBA
        //color.RGBA{n(), n(), n(), 0xff}
	My_go_color color.RGBA
        //Slice(property string, result interface{})
	//MyGoSlice
	//other golang object structures
}

func (v *GoType) SetMy_go_text(someText string) {
	fmt.Println("Text changing to:", someText)
	v.My_go_text = someText
}

func (v *GoType) SetMy_go_int(someInt_ int) {
	fmt.Println("My_go_int changing to:", someInt_)
	v.My_go_int = someInt_
}

func (v *GoType) SetMy_go_int64(someInt64_ int64) {
	fmt.Println("My_go_int64 changing to:", someInt64_)
	v.My_go_int64 = someInt64_
}

func (v *GoType) SetMy_go_float64(someFloat64_ float64) {
	fmt.Println("My_go_float64 changing to:", someFloat64_)
	v.My_go_float64 = someFloat64_
}

func (v *GoType) SetMy_go_bool(someBool_ bool) {
	fmt.Println("My_go_bool changing to:", someBool_)
	v.My_go_bool = someBool_
}

func (v *GoType) SetMy_go_string(someString_ string) {
	fmt.Println("My_go_string changing to:", someString_)
	v.My_go_string = someString_
}

func (v *GoType) SetMy_go_color(someColor_ color.RGBA) {
	fmt.Println("My_go_color changing to:", someColor_)
	v.My_go_color = someColor_
}

//add some other type setters here
//Slice(property string, result interface{})
//MyGoSlice
//other golang object structures

type GoSingleton struct {
	Event string
}


//this sets up the qml engine
//registers some golang types into the qml engine
//then loads up a qml component file
//Once the qml component file is loaded,
//we create a qml component object instance,
//then we check to see the golang bridge types' values
//before exiting the program.
//NOTE:  It seems the actual golang type is instantiated on the
//qml engine side which is very cool.
func run() error {
	qml.Init(nil)

	qml.RegisterTypes("GoExtensions", 1, 0, []qml.TypeSpec{{
		Init: func(v *GoType, obj qml.Object) {},
	}, {
		Init: func(v *GoSingleton, obj qml.Object) { v.Event = "birthday" },

		Singleton: true,
	}})

	engine := qml.NewEngine()
	component, err := engine.LoadFile("customtype.qml")
	if err != nil {
		return err
	}

	value := component.Create(nil)
	fmt.Println("My_go_text is:", value.Interface().(*GoType).My_go_text)
	fmt.Println("My_go_int is:", value.Interface().(*GoType).My_go_int)
	fmt.Println("My_go_int64 is:", value.Interface().(*GoType).My_go_int64)
	fmt.Println("My_go_float64 is:", value.Interface().(*GoType).My_go_float64)
	fmt.Println("My_go_bool is:", value.Interface().(*GoType).My_go_bool)
	fmt.Println("My_go_string is:", value.Interface().(*GoType).My_go_string)
	fmt.Println("My_go_color is:", value.Interface().(*GoType).My_go_color)

	return nil
}

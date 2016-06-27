package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"time"

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
	colors := NewList(engine)
	engine.Context().SetVar("colors", colors)
	cls := engine.Context().Var("colors")
	fmt.Printf("Colors: %T %v\n", cls, cls)
	component, err := engine.LoadFile("delegate.qml")
	if err != nil {
		return err
	}
	window := component.CreateWindow(nil)
	window.Show()
	go func() {
		n := func() uint8 { return uint8(rand.Intn(256)) }
		for i := 0; i < 100; i++ {
			colors.Add(color.RGBA{n(), n(), n(), 0xff})
			time.Sleep(1 * time.Second)
		}
	}()
	window.Wait()
	return nil
}

type List struct {
	qml.ItemModel
	internal qml.ItemModelInternal
	list     []color.RGBA
}

var _ qml.ItemModelImpl = &List{}

func NewList(engine *qml.Engine) *List {
	list := &List{}
	list.ItemModel, list.internal = qml.NewItemModel(engine, nil, list)
	return list
}

func (l *List) Add(c color.RGBA) {
	qml.RunMain(func() {
		l.internal.BeginInsertRows(nil, len(l.list), len(l.list))
		l.list = append(l.list, c)
		l.internal.EndInsertRows()
	})
}

func (l *List) validateModelIndex(mi qml.ModelIndex) bool {
	column := mi.Column()
	row := mi.Row()
	return mi.IsValid() && column == 0 && row >= 0 && row < len(l.list)
}

func fmtMI(mi qml.ModelIndex) string {
	if !mi.IsValid() {
		return "ModelIndex{invalid}"
	}
	return fmt.Sprintf("ModelIndex{%v, %v, %v}", mi.Row(), mi.Column(), mi.InternalId())
}

// Required functions
func (l *List) ColumnCount(parent qml.ModelIndex) int {
	fmt.Println("ColumnCount:", fmtMI(parent))
	return 1
}

func (l *List) RowCount(parent qml.ModelIndex) int {
	fmt.Println("RowCount:", fmtMI(parent), len(l.list))
	return len(l.list)
}

func (l *List) Data(index qml.ModelIndex, role qml.Role) interface{} {
	fmt.Println("Data:", fmtMI(index), role)
	return l.list[index.Row()]
}

func (l *List) Index(row int, column int, parent qml.ModelIndex) qml.ModelIndex {
	// fmt.Println("Index:", row, column, fmtMI(parent))
	if !parent.IsValid() && column == 0 && row >= 0 && row < len(l.list) {
		return l.internal.CreateIndex(row, column, 0)
	}
	return nil
}

func (l *List) Parent(index qml.ModelIndex) qml.ModelIndex {
	// fmt.Println("Parent:", index)
	return nil
}

// Required for editing
func (l *List) Flags(index qml.ModelIndex) qml.ItemFlags {
	fmt.Println("Flags:", fmtMI(index))
	return qml.NoItemFlags
}

func (l *List) SetData(index qml.ModelIndex, value interface{}, role qml.Role) bool {
	fmt.Println("SetData:", fmtMI(index), role, value)
	if !l.validateModelIndex(index) {
		return false
	}
	rgba, ok := value.(color.RGBA)
	if !ok {
		return false
	}
	l.list[index.Row()] = rgba
	return true
}

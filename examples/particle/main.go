package main

import (
	"gopkg.in/qml.v0"
	"math/rand"
	"time"
)

func main() {
	qml.Init(nil)
	engine := qml.NewEngine()
	component, err := engine.LoadFile("particle.qml")
	if err != nil {
		panic(err)
	}

	ctrl := Control{Message: "Hello from Go!"}

	context := engine.Context()
	context.SetVar("ctrl", &ctrl)

	window := component.CreateWindow(nil)

	ctrl.Root = window.Root()

	rand.Seed(time.Now().Unix())

	window.Show()
	window.Wait()
}

type Control struct {
	Root    qml.Object
	Message string
}

func (ctrl *Control) TextReleased(text qml.Object) {
	x := text.Int("x")
	y := text.Int("y")
	width := text.Int("width")
	height := text.Int("height")

	ctrl.Emit(x+15, y+height/2)
	ctrl.Emit(x+width/2, 1.0*y+height/2)
	ctrl.Emit(x+width-15, 1.0*y+height/2)

	go func() {
		time.Sleep(500 * time.Millisecond)
		messages := []string{"Hello", "Hello", "Hacks"}
		ctrl.Message = messages[rand.Intn(len(messages))] + " from Go!"
		qml.Changed(ctrl, &ctrl.Message)
	}()
}

func (ctrl *Control) Emit(x, y int) {
	component := ctrl.Root.Object("emitterComponent")
	for i := 0; i < 8; i++ {
		emitter := component.Create(nil)
		emitter.Set("x", x)
		emitter.Set("y", y)
		emitter.Set("targetX", rand.Intn(240)-120+x)
		emitter.Set("targetY", rand.Intn(240)-120+y)
		emitter.Set("life", rand.Intn(2400)+200)
		emitter.Set("emitRate", rand.Intn(32)+32)
		emitter.ObjectByName("xAnim").Call("start")
		emitter.ObjectByName("yAnim").Call("start")
		emitter.Set("enabled", true)
	}
}

func (ctrl *Control) Done(emitter qml.Object) {
	emitter.Destroy()
}

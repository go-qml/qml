package main

import (
	"fmt"
	"github.com/niemeyer/qml"
	"github.com/go-gl/gl"
	"os"
)

func main() {
	filename := "gopher.qml"
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}
	if err := run(filename); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(filename string) error {
	qml.Init(nil)
	engine := qml.NewEngine()

	model, err := Read("model/gopher.obj")
	if err != nil {
		return err
	}

	qml.RegisterTypes("GoExtensions", 1, 0, []qml.TypeSpec{{
		Init: func(g *Gopher, obj qml.Object) {
			g.Object = obj
			g.model = model
		},
	}})

	component, err := engine.LoadFile(filename)
	if err != nil {
		return err
	}

	win := component.CreateWindow(nil)
	win.Set("x", 560)
	win.Set("y", 320)
	win.Show()
	win.Wait()
	return nil
}

type Gopher struct {
	qml.Object

	model map[string]*Object

	Rotation int
}

func (r *Gopher) SetRotation(rotation int) {
	r.Rotation = rotation
	r.Call("update")
}

func (r *Gopher) Paint(p *qml.Painter) {
	width := float32(r.Int("width"))

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.ShadeModel(gl.SMOOTH)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(true)
	gl.Enable(gl.NORMALIZE)

	gl.Clear(gl.DEPTH_BUFFER_BIT)

	gl.Scalef(width/3, width/3, width/3)

	lka := []float32{0.3, 0.3, 0.3, 1.0}
	lkd := []float32{1.0, 1.0, 1.0, 0.0}
	lks := []float32{1.0, 1.0, 1.0, 1.0}
	lpos := []float32{-2, 6, 3, 1.0}

	gl.Enable(gl.LIGHTING)
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, lka)
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, lkd)
	gl.Lightfv(gl.LIGHT0, gl.SPECULAR, lks)
	gl.Lightfv(gl.LIGHT0, gl.POSITION, lpos)
	gl.Enable(gl.LIGHT0)

	gl.EnableClientState(gl.NORMAL_ARRAY)
	gl.EnableClientState(gl.VERTEX_ARRAY)

	gl.Translatef(1.5, 1.5, 0)
	gl.Rotatef(-90, 0, 0, 1)
	gl.Rotatef(float32(90+((36000+r.Rotation)%360)), 1, 0, 0)

	gl.Disable(gl.COLOR_MATERIAL)

	for _, obj := range r.model {
		for _, group := range obj.Groups {
			gl.Materialfv(gl.FRONT, gl.AMBIENT, group.Material.Ambient)
			gl.Materialfv(gl.FRONT, gl.DIFFUSE, group.Material.Diffuse)
			gl.Materialfv(gl.FRONT, gl.SPECULAR, group.Material.Specular)
			gl.Materialf(gl.FRONT, gl.SHININESS, group.Material.Shininess)
			gl.VertexPointer(3, gl.FLOAT, 0, group.Vertexes)
			gl.NormalPointer(gl.FLOAT, 0, group.Normals)
			gl.DrawArrays(gl.TRIANGLES, 0, len(group.Vertexes)/3)
		}
	}

	gl.Enable(gl.COLOR_MATERIAL)

	gl.DisableClientState(gl.NORMAL_ARRAY)
	gl.DisableClientState(gl.VERTEX_ARRAY)
}

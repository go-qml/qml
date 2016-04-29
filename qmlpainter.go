package qml

import "github.com/limetext/qml-go/gl/glbase"

// Painter is provided to Paint methods on Go types that have displayable content.
type Painter struct {
	engine *Engine
	obj    Object
	glctxt glbase.Context
}

// Object returns the underlying object being painted.
func (p *Painter) Object() Object {
	return p.obj
}

// GLContext returns the OpenGL context for this painter.
func (p *Painter) GLContext() *glbase.Context {
	return &p.glctxt
}

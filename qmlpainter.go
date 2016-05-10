package qml

// #include <stdlib.h>
//
// #include "capi.h"
//
// #include "painter.h"
//
import "C"

import (
	"unsafe"

	"github.com/limetext/qml-go/gl/glbase"
)

type RenderHint C.enum_RenderHint

const /*RenderHint*/ (
	Antialiasing            RenderHint = C.Antialiasing
	TextAntialiasing        RenderHint = C.TextAntialiasing
	SmoothPixmapTransform   RenderHint = C.SmoothPixmapTransform
	HighQualityAntialiasing RenderHint = C.HighQualityAntialiasing
	NonCosmeticDefaultPen   RenderHint = C.NonCosmeticDefaultPen
	Qt4CompatiblePainting   RenderHint = C.Qt4CompatiblePainting
)

type CompositionMode C.enum_CompositionMode

const /* CompositionMode */ (
	CompositionMode_SourceOver      CompositionMode = C.CompositionMode_SourceOver
	CompositionMode_DestinationOver CompositionMode = C.CompositionMode_DestinationOver
	CompositionMode_Clear           CompositionMode = C.CompositionMode_Clear
	CompositionMode_Source          CompositionMode = C.CompositionMode_Source
	CompositionMode_Destination     CompositionMode = C.CompositionMode_Destination
	CompositionMode_SourceIn        CompositionMode = C.CompositionMode_SourceIn
	CompositionMode_DestinationIn   CompositionMode = C.CompositionMode_DestinationIn
	CompositionMode_SourceOut       CompositionMode = C.CompositionMode_SourceOut
	CompositionMode_DestinationOut  CompositionMode = C.CompositionMode_DestinationOut
	CompositionMode_SourceAtop      CompositionMode = C.CompositionMode_SourceAtop
	CompositionMode_DestinationAtop CompositionMode = C.CompositionMode_DestinationAtop
	CompositionMode_Xor             CompositionMode = C.CompositionMode_Xor

	//svg 1.2 blend modes
	CompositionMode_Plus       CompositionMode = C.CompositionMode_Plus
	CompositionMode_Multiply   CompositionMode = C.CompositionMode_Multiply
	CompositionMode_Screen     CompositionMode = C.CompositionMode_Screen
	CompositionMode_Overlay    CompositionMode = C.CompositionMode_Overlay
	CompositionMode_Darken     CompositionMode = C.CompositionMode_Darken
	CompositionMode_Lighten    CompositionMode = C.CompositionMode_Lighten
	CompositionMode_ColorDodge CompositionMode = C.CompositionMode_ColorDodge
	CompositionMode_ColorBurn  CompositionMode = C.CompositionMode_ColorBurn
	CompositionMode_HardLight  CompositionMode = C.CompositionMode_HardLight
	CompositionMode_SoftLight  CompositionMode = C.CompositionMode_SoftLight
	CompositionMode_Difference CompositionMode = C.CompositionMode_Difference
	CompositionMode_Exclusion  CompositionMode = C.CompositionMode_Exclusion

	// ROPs
	RasterOp_SourceOrDestination        CompositionMode = C.RasterOp_SourceOrDestination
	RasterOp_SourceAndDestination       CompositionMode = C.RasterOp_SourceAndDestination
	RasterOp_SourceXorDestination       CompositionMode = C.RasterOp_SourceXorDestination
	RasterOp_NotSourceAndNotDestination CompositionMode = C.RasterOp_NotSourceAndNotDestination
	RasterOp_NotSourceOrNotDestination  CompositionMode = C.RasterOp_NotSourceOrNotDestination
	RasterOp_NotSourceXorDestination    CompositionMode = C.RasterOp_NotSourceXorDestination
	RasterOp_NotSource                  CompositionMode = C.RasterOp_NotSource
	RasterOp_NotSourceAndDestination    CompositionMode = C.RasterOp_NotSourceAndDestination
	RasterOp_SourceAndNotDestination    CompositionMode = C.RasterOp_SourceAndNotDestination
	RasterOp_NotSourceOrDestination     CompositionMode = C.RasterOp_NotSourceOrDestination
	RasterOp_SourceOrNotDestination     CompositionMode = C.RasterOp_SourceOrNotDestination
	RasterOp_ClearDestination           CompositionMode = C.RasterOp_ClearDestination
	RasterOp_SetDestination             CompositionMode = C.RasterOp_SetDestination
	RasterOp_NotDestination             CompositionMode = C.RasterOp_NotDestination
)

// Painter is provided to Paint methods on Go types that have displayable content.
type Painter struct {
	engine   *Engine
	obj      Object
	glctxt   glbase.Context
	qpainter unsafe.Pointer
}

// Object returns the underlying object being painted.
func (p *Painter) Object() Object {
	return p.obj
}

// GLContext returns the OpenGL context for this painter.
func (p *Painter) GLContext() *glbase.Context {
	return &p.glctxt
}

func (p *Painter) CompositionMode() CompositionMode {
	return CompositionMode(C.painterCompositionMode(p.qpainter))
}

func (p *Painter) SetCompositionMode(mode CompositionMode) {
	C.painterSetCompositionMode(p.qpainter, C.enum_CompositionMode(mode))
}

func (p *Painter) Save() {
	C.painterSave(p.qpainter)
}
func (p *Painter) Restore() {
	C.painterRestore(p.qpainter)
}

func (p *Painter) BeginNativePainting() {
	C.painterBeginNativePainting(p.qpainter)
}

func (p *Painter) EndNativePainting() {
	C.painterEndNativePainting(p.qpainter)
}

package qimage

// #cgo CPPFLAGS: -I../cpp
// #cgo CXXFLAGS: -std=c++0x -pedantic-errors -Wall -fno-strict-aliasing
// #cgo LDFLAGS: -lstdc++
// #cgo pkg-config: Qt5Core Qt5Widgets Qt5Quick
//
// #include <stdlib.h>
// #include "image.h"
//
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/limetext/qml-go/internal/util"
)

type QImage interface {
}

type Format uint

const (
	Format_Invalid              Format = 0
	Format_RGB32                Format = 4
	Format_ARGB32_Premultiplied Format = 6
	Format_RGB16                Format = 7
)

type qImage struct {
	ptr unsafe.Pointer
}

func NewImage(width int, height int, format Format) QImage {
	qimgptr := C.newQImage(C.int(width), C.int(height), C.uint(format))

	return mkImage(qimgptr)
}

func LoadImage(filename string, format string) QImage {
	fname, fnamelen := util.UnsafeStringData(filename)
	fmat, _ := util.UnsafeStringData(format + "\x00")

	qimgptr := C.loadQImage((*C.char)(fname), C.int(fnamelen), (*C.char)(fmat))

	return mkImage(qimgptr)
}

func mkImage(qimgptr unsafe.Pointer) *qImage {
	qimg := &qImage{ptr: qimgptr}

	runtime.SetFinalizer(qimg, deleteImage)

	return qimg
}

func deleteImage(qimg *qImage) {
	C.deleteQImage(qimg.ptr)
	qimg.ptr = nil
}

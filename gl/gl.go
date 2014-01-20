package gl

// #cgo CXXFLAGS: -std=c++0x -pedantic-errors -Wall -fno-strict-aliasing -DGL_GLEXT_PROTOTYPES
// #cgo LDFLAGS: -lstdc++ -lGL
// #cgo pkg-config: Qt5Core Qt5Widgets Qt5Quick Qt5OpenGL
//
// #define GL_GLEXT_PROTOTYPES
// #include "gl.h"
// #include "glext.h"
//
import "C"

import (
	"fmt"
	"unsafe"
)

func GenBuffers(n Sizei, buffers []Uint) {
	if Sizei(len(buffers)) < n {
		panic(fmt.Sprintf("trying to use %d elements in a buffer of length %d", n, len(buffers)))
	}
	C.glGenBuffers(C.GLsizei(n), (*C.GLuint)(unsafe.Pointer(&buffers[0])))
}

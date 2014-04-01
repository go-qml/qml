package gl

// #cgo CXXFLAGS: -std=c++0x -pedantic-errors -Wall -fno-strict-aliasing -DGL_GLEXT_PROTOTYPES -DGL_GLEXT_LEGACY
// #cgo LDFLAGS: -lstdc++
// #cgo !darwin LDFLAGS: -lGL
// #cgo  darwin LDFLAGS: -framework OpenGL
// #cgo pkg-config: Qt5Core Qt5Widgets Qt5Quick Qt5OpenGL
//
// #include "gl.h"
//
import "C"

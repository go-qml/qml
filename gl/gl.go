package gl

// #cgo CXXFLAGS: -std=c++0x -pedantic-errors -Wall -fno-strict-aliasing 
// #cgo CPPFLAGS: -DGL_GLEXT_PROTOTYPES -DGL_GLEXT_LEGACY
// #cgo LDFLAGS: -lstdc++
// #cgo pkg-config: Qt5Core Qt5Widgets Qt5Quick Qt5OpenGL
//
// #include "gl.h"
//
import "C"

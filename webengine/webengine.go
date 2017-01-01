package webengine

// #cgo CPPFLAGS: -I./
// #cgo CXXFLAGS: -std=c++0x -pedantic-errors -Wall -fno-strict-aliasing
// #cgo LDFLAGS: -lstdc++
// #cgo pkg-config: Qt5WebEngine
//
// #include "webengine.h"
import "C"

import (
	"gopkg.in/qml.v1"
)

// Initializes the WebEngine extension.
func Initialize() {
	qml.RunMain(func() {
		C.webengineInitialize()
	})
}

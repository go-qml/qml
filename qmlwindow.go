package qml

// #include <stdlib.h>
//
// #include "capi.h"
//
import "C"

import (
	"image"
	"reflect"
	"sync"
	"unsafe"
)

// Window represents a QML window where components are rendered.
type Window struct {
	Common
}

// Show exposes the window.
func (win *Window) Show() {
	RunMain(func() {
		C.windowShow(win.addr)
	})
}

// Hide hides the window.
func (win *Window) Hide() {
	RunMain(func() {
		C.windowHide(win.addr)
	})
}

// PlatformId returns the window's platform id.
//
// For platforms where this id might be useful, the value returned will
// uniquely represent the window inside the corresponding screen.
func (win *Window) PlatformId() uintptr {
	var id uintptr
	RunMain(func() {
		id = uintptr(C.windowPlatformId(win.addr))
	})
	return id
}

// Root returns the root object being rendered.
//
// If the window was defined in QML code, the root object is the window itself.
func (win *Window) Root() Object {
	var obj Common
	obj.engine = win.engine
	RunMain(func() {
		obj.setAddr(C.windowRootObject(win.addr))
	})
	return &obj
}

// Wait blocks the current goroutine until the window is closed.
func (win *Window) Wait() {
	// XXX Test this.
	var m sync.Mutex
	m.Lock()
	RunMain(func() {
		// TODO Must be able to wait for the same Window from multiple goroutines.
		// TODO If the window is not visible, must return immediately.
		waitingWindows[win.addr] = &m
		C.windowConnectHidden(win.addr)
	})
	m.Lock()
}

var waitingWindows = make(map[unsafe.Pointer]*sync.Mutex)

//export hookWindowHidden
func hookWindowHidden(addr unsafe.Pointer) {
	m, ok := waitingWindows[addr]
	if !ok {
		panic("window is not waiting")
	}
	delete(waitingWindows, addr)
	m.Unlock()
}

// Snapshot returns an image with the visible contents of the window.
// The main GUI thread is paused while the data is being acquired.
func (win *Window) Snapshot() image.Image {
	// TODO Test this.
	var cimage unsafe.Pointer
	RunMain(func() {
		cimage = C.windowGrabWindow(win.addr)
	})
	defer C.delImage(cimage)

	// This should be safe to be done out of the main GUI thread.
	var cwidth, cheight C.int
	C.imageSize(cimage, &cwidth, &cheight)

	var cbits []byte
	cbitsh := (*reflect.SliceHeader)((unsafe.Pointer)(&cbits))
	cbitsh.Data = (uintptr)((unsafe.Pointer)(C.imageConstBits(cimage)))
	cbitsh.Len = int(cwidth * cheight * 8) // ARGB
	cbitsh.Cap = cbitsh.Len

	image := image.NewRGBA(image.Rect(0, 0, int(cwidth), int(cheight)))
	l := int(cwidth * cheight * 4)
	for i := 0; i < l; i += 4 {
		var c uint32 = *(*uint32)(unsafe.Pointer(&cbits[i]))
		image.Pix[i+0] = byte(c >> 16)
		image.Pix[i+1] = byte(c >> 8)
		image.Pix[i+2] = byte(c)
		image.Pix[i+3] = byte(c >> 24)
	}
	return image
}

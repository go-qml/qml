package qml

// #cgo CPPFLAGS: -I/usr/include/qt5/QtCore/5.0.2/QtCore
// #cgo pkg-config: Qt5Core Qt5Widgets Qt5Quick
// #cgo LDFLAGS: -lstdc++
//
// #include "capi.h"
//
import "C"

import (
	"reflect"
	"runtime"
	"unsafe"
)


// guiLoop runs the main GUI thread event loop in C++ land.
func guiLoop() {
	runtime.LockOSThread()
	C.newGuiApplication()
	C.startIdleTimer()
	C.applicationExec();
}

var guiLock = 0
var guiFunc = make(chan func())
var guiDone = make(chan struct{}, 1)

// gui runs f in the main GUI thread and waits for f to return.
func gui(f func()) {
	guiFunc <- f
	<-guiDone
}

// Lock freezes all QML activity by blocking the main event loop.
// Locking is necessary before updating shared data structures
// without race conditions.
//
// It's safe to use qml functionality while holding a lock, as
// long as the requests made to not depend on follow up events
// to be processed.
//
// The Lock method is reentrant. That is, the method may be called
// multiple times, and QML activities will only be resumed after
// Unlock is called a matching number of times.
func Lock() {
	gui(func() {
		guiLock++
	})
}

// Unlock releases the QML event loop. See Lock for details.
func Unlock() {
	gui(func() {
		if guiLock == 0 {
			panic("qml.Unlock called without lock held")
		}
		guiLock--
	})
}

// hookIdleTimer is run once per iteration of the Qt event loop, within
// the main GUI thread.
//
//export hookIdleTimer
func hookIdleTimer() {
	var f func()
	for {
		select {
		case f = <-guiFunc:
		default:
			if guiLock > 0 {
				f = <-guiFunc
			} else {
				return
			}
		}
		f()
		guiDone <- struct{}{}
	}
}

//export hookReadField
func hookReadField(ifacep unsafe.Pointer, memberIndex C.int, result *C.DataValue) {
	value := *(*interface{})(ifacep)
	//fmt.Printf("QML requested member %d for Go's %T at %p.\n", memberIndex, *ifacep, ifacep)
	field := reflect.ValueOf(value).Elem().Field(int(memberIndex))

	// TODO Strings are being passed in an unsafe manner here. There is a
	// small chance that the field is changed and the garbage collector run
	// before C++ has chance to look at the data. We can solve this problem
	// by queuing up values in a stack, and cleaning the stack when the
	// idle timer fires next.
	packDataValue(field.Interface(), result)
}

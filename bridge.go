package qml

// #cgo CPPFLAGS: -I/usr/include/qt5/QtCore/5.0.2/QtCore
// #cgo CXXFLAGS: -std=c++11
// #cgo pkg-config: Qt5Core Qt5Widgets Qt5Quick glib-2.0
// #cgo LDFLAGS: -lstdc++
//
// #include "capi.h"
//
import "C"

import (
	"fmt"
	"reflect"
	"runtime"
	"sync/atomic"
	"unsafe"
)

var hookWaiting C.int

// guiLoop runs the main GUI thread event loop in C++ land.
func guiLoop() {
	runtime.LockOSThread()
	C.newGuiApplication()
	C.startIdleTimer(&hookWaiting)
	C.applicationExec()
}

var guiFunc = make(chan func())
var guiDone = make(chan struct{})
var guiLock = 0

// gui runs f in the main GUI thread and waits for f to return.
func gui(f func()) {
	// Tell Qt we're waiting for the idle hook to be called.
	atomic.AddInt32((*int32)(unsafe.Pointer(&hookWaiting)), 1)

	// Send f to be executed by the idle hook in the main GUI thread.
	guiFunc <- f

	// Wait until f is done executing.
	<-guiDone
}

// Lock freezes all QML activity by blocking the main event loop.
// Locking is necessary before updating shared data structures
// without race conditions.
//
// It's safe to use qml functionality while holding a lock, as
// long as the requests made do not depend on follow up QML
// events to be processed before returning. If that happens, the
// problem will be observed as the application freezing.
//
// The Lock function is reentrant. That means it may be called
// multiple times, and QML activities will only be resumed after
// Unlock is called a matching number of times.
func Lock() {
	// TODO Better testing for this.
	gui(func() {
		guiLock++
	})
}

// Unlock releases the QML event loop. See Lock for details.
func Unlock() {
	gui(func() {
		if guiLock == 0 {
			panic("qml.Unlock called without lock being held")
		}
		guiLock--
	})
}

// FlushAll synchronously flushes all pending QML activities.
func FlushAll() {
	// TODO Better testing for this.
	gui(func() {
		C.applicationFlushAll()
	})
}

// hookIdleTimer is run once per iteration of the Qt event loop,
// within the main GUI thread, but only if at least one goroutine
// has atomically incremented hookWaiting.
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
		atomic.AddInt32((*int32)(unsafe.Pointer(&hookWaiting)), -1)
	}
}

type cppGoValue struct {
	ifacep *interface{}
	valuep unsafe.Pointer
	owner  valueOwner
}

type valueOwner int

const (
	anyOwner = iota
	cppOwner
	jsOwner
)

// wrapGoValue creates a new GoValue object in C++ land wrapping
// the Go value contained in the given interface.
//
// This must be run from the main GUI thread.
func wrapGoValue(engine *Engine, value interface{}, owner valueOwner) (valuep unsafe.Pointer) {
	gv, ok := engine.values[value]
	if !ok {
		parent := nilPtr
		if owner == cppOwner {
			parent = engine.addr
		}
		// Define a local copy rather than using &value directly, to
		// avoid forcing value's off-stack allocation unnecessarily.
		iface := value
		gv.ifacep = &iface
		gv.valuep = C.newValue(unsafe.Pointer(&iface), typeInfo(value), parent)
		gv.owner = owner
		engine.values[value] = gv
		stats.valuesAlive(+1)
		C.engineSetContextForObject(engine.addr, gv.valuep);
		switch owner {
		case cppOwner:
			C.engineSetOwnershipCPP(engine.addr, gv.valuep)
		case jsOwner:
			C.engineSetOwnershipJS(engine.addr, gv.valuep)
		}
	} else if owner == cppOwner && gv.owner != cppOwner {
		gv.owner = cppOwner
		C.engineSetOwnershipCPP(engine.addr, gv.valuep)
		C.objectSetParent(gv.valuep, engine.addr)
	}
	return gv.valuep
}

//export hookGoValueDestroyed
func hookGoValueDestroyed(enginep unsafe.Pointer, ifacep unsafe.Pointer) {
	fmt.Println("GoValue destroyed!")
	engine := engines[enginep]
	value := *(*interface{})(ifacep)
	if engine == nil {
		panic("unknown engine pointer; who created it?")
	}
	before := len(engine.values)
	delete(engine.values, value)
	if len(engine.values) == before {
		panic("unknown value was destroyed")
	}
	stats.valuesAlive(-1)
}

//export hookGoValueReadField
func hookGoValueReadField(enginep unsafe.Pointer, ifacep unsafe.Pointer, memberIndex C.int, result *C.DataValue) {
	engine := engines[enginep]
	value := *(*interface{})(ifacep)

	if engine == nil {
		if enginep == nilPtr {
			panic("nil engine pointer; who created the object?")
		} else {
			panic("unknown engine pointer; who created the engine?")
		}
	}

	v := reflect.ValueOf(value)
	for v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	field := v.Field(int(memberIndex))

	// TODO Strings are being passed in an unsafe manner here. There is a
	// small chance that the field is changed and the garbage collector is run
	// before C++ has a chance to look at the data. We can solve this problem
	// by queuing up values in a stack, and cleaning the stack when the
	// idle timer fires next.
	packDataValue(field.Interface(), result, engine, jsOwner)
}

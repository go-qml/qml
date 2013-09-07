package qml

// #cgo CPPFLAGS: -I/usr/include/qt5/QtCore/5.0.2/QtCore -I./cpp
// #cgo CXXFLAGS: -std=c++11
// #cgo pkg-config: Qt5Core Qt5Widgets Qt5Quick glib-2.0
// #cgo LDFLAGS: -lstdc++
//
// #include "cpp/capi.h"
//
import "C"

import (
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

func Notify(value interface{}, field string) {

	// TODO Must notify all engines, not one of them.

	// TODO Must access engine.values in a non-racy way.
	var fold *valueFold
	for _, engine := range engines {
		if fold = engine.values[value]; fold != nil {
			break
		}
	}

	if fold == nil {

		for f, _ := range enginePending {
			if f.gvalue == value {
				fold = f
				break
			}
		}

		if fold == nil {
			// TODO Perhaps return an error instead.
			panic("value is not known")
		}
	}

	// TODO Can probably use the field address for notify, as in:
	//          Notify(&value, &value.field)
	//      And do it in O(1).
	vt := reflect.ValueOf(value).Type()
	for vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}
	numField := vt.NumField()
	for i := 0; i < numField; i++ {
		if vt.Field(i).Name == field {
			gui(func() {
				C.goValueActivate(fold.cvalue, C.int(i))
			})
		}
	}
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

type valueFold struct {
	engine *Engine
	gvalue interface{}
	cvalue unsafe.Pointer
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
func wrapGoValue(engine *Engine, gvalue interface{}, owner valueOwner) (cvalue unsafe.Pointer) {
	// TODO Return an error if gvalue is a non-basic type and not a pointer.
	//      Pointer-to-pointer is also not okay.
	fold, ok := engine.values[gvalue]
	if !ok {
		parent := nilPtr
		if owner == cppOwner {
			parent = engine.addr
		}
		fold = &valueFold{
			engine: engine,
			gvalue: gvalue,
			owner:  owner,
		}
		fold.cvalue = C.newGoValue(unsafe.Pointer(fold), typeInfo(gvalue), parent)
		engine.values[gvalue] = fold
		stats.valuesAlive(+1)
		C.engineSetContextForObject(engine.addr, fold.cvalue)
		switch owner {
		case cppOwner:
			C.engineSetOwnershipCPP(engine.addr, fold.cvalue)
		case jsOwner:
			C.engineSetOwnershipJS(engine.addr, fold.cvalue)
		}
	} else if owner == cppOwner && fold.owner != cppOwner {
		fold.owner = cppOwner
		C.engineSetOwnershipCPP(engine.addr, fold.cvalue)
		C.objectSetParent(fold.cvalue, engine.addr)
	}
	return fold.cvalue
}

var enginePending = make(map[*valueFold]bool)

//export hookGoValueTypeNew
func hookGoValueTypeNew(cvalue unsafe.Pointer, specp unsafe.Pointer) (foldp unsafe.Pointer) {
	fold := &valueFold{
		gvalue: (*TypeSpec)(specp).New(),
		cvalue: cvalue,
		owner:  jsOwner,
	}
	enginePending[fold] = true
	stats.valuesAlive(+1)
	return unsafe.Pointer(fold)
}

//export hookGoValueDestroyed
func hookGoValueDestroyed(enginep unsafe.Pointer, foldp unsafe.Pointer) {
	fold := (*valueFold)(foldp)
	engine := fold.engine
	if engine == nil {
		before := len(enginePending)
		delete(enginePending, fold)
		if len(enginePending) == before {
			panic("destroying value without an associated engine and unknown to the pending engine set; who created the value?")
		}
	} else if engines[engine.addr] == nil {
		// Must never do that. The engine holds memory references that C++ depends on.
		panic("engine was released from global list while its values were still alive")
	} else {
		before := len(engine.values)
		delete(engine.values, fold.gvalue)
		if len(engine.values) == before {
			panic("destroying value that knows about the engine, but the engine doesn't know about the value; who cleared the engine?")
		}
		if engine.destroyed && len(engine.values) == 0 {
			delete(engines, engine.addr)
		}
	}
	stats.valuesAlive(-1)
}

//export hookGoValueReadField
func hookGoValueReadField(enginep unsafe.Pointer, foldp unsafe.Pointer, reflectIndex C.int, resultdv *C.DataValue) {
	fold := ensureEngine(enginep, foldp)
	v := reflect.ValueOf(fold.gvalue)
	for v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	field := v.Field(int(reflectIndex))

	// TODO Strings are being passed in an unsafe manner here. There is a
	// small chance that the field is changed and the garbage collector is run
	// before C++ has a chance to look at the data. We can solve this problem
	// by queuing up values in a stack, and cleaning the stack when the
	// idle timer fires next.
	packDataValue(field.Interface(), resultdv, fold.engine, jsOwner)
}

//export hookGoValueWriteField
func hookGoValueWriteField(enginep unsafe.Pointer, foldp unsafe.Pointer, reflectIndex C.int, assigndv *C.DataValue) {
	fold := ensureEngine(enginep, foldp)
	v := reflect.ValueOf(fold.gvalue)
	for v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	field := v.Field(int(reflectIndex))
	assign := unpackDataValue(assigndv)

	// TODO What to do if it fails?
	convertAndSet(field, reflect.ValueOf(assign))
}

func convertAndSet(to, from reflect.Value) {
	// TODO Catch the panic and error out.
	to.Set(from.Convert(to.Type()))
}

//export hookGoValueCallMethod
func hookGoValueCallMethod(enginep unsafe.Pointer, foldp unsafe.Pointer, reflectIndex C.int, resultdv *C.DataValue) {
	fold := ensureEngine(enginep, foldp)
	v := reflect.ValueOf(fold.gvalue)

	// TODO Must ensure that v is necessarily a pointer here.

	method := v.Method(int(reflectIndex))

	// TODO Unhardcode this.
	result := method.Call(nil)
	if len(result) != 1 || result[0].Type() != typeString {
		panic("result must be a string for now")
	}

	packDataValue(result[0].Interface(), resultdv, fold.engine, jsOwner)
}

func ensureEngine(enginep unsafe.Pointer, foldp unsafe.Pointer) *valueFold {
	fold := (*valueFold)(foldp)
	if fold.engine == nil {
		if enginep == nilPtr {
			panic("accessing field from value without an engine pointer; who created the value?")
		}
		engine := engines[enginep]
		if engine == nil {
			panic("unknown engine pointer; who created the engine?")
		}
		fold.engine = engine
		engine.values[fold.gvalue] = fold
		before := len(enginePending)
		delete(enginePending, fold)
		if len(enginePending) == before {
			panic("value had no engine, but is not in the pending engine set; who created the value?")
		}
	}
	return fold
}

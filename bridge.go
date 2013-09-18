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
	"fmt"
	"github.com/niemeyer/qml/tref"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

var hookWaiting C.int

// guiLoop runs the main GUI thread event loop in C++ land.
func guiLoop() {
	runtime.LockOSThread()
	guiLoopRef = tref.Ref()
	guiLoopReady.Unlock()
	C.newGuiApplication()
	C.startIdleTimer(&hookWaiting)
	C.applicationExec()
}

var (
	guiFunc      = make(chan func())
	guiDone      = make(chan struct{})
	guiLock      = 0
	guiLoopReady sync.Mutex
	guiLoopRef   uintptr
)

// gui runs f in the main GUI thread and waits for f to return.
func gui(f func()) {
	if tref.Ref() == guiLoopRef {
		// Already within the GUI thread. Attempting to wait would deadlock.
		f()
		return
	}

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

// Flush synchronously flushes all pending QML activities.
func Flush() {
	// TODO Better testing for this.
	gui(func() {
		C.applicationFlushAll()
	})
}

// Changed notifies all QML bindings that the given field value has changed.
//
// For example:
//
//     qml.Changed(&value, &value.Field)
//
func Changed(value, fieldAddr interface{}) {
	valuev := reflect.ValueOf(value)
	fieldv := reflect.ValueOf(fieldAddr)
	for valuev.Kind() == reflect.Ptr {
		valuev = valuev.Elem()
	}
	for fieldv.Kind() == reflect.Ptr {
		fieldv = fieldv.Elem()
	}
	if fieldv.Type().Size() == 0 {
		panic("cannot report changes on zero-sized fields")
	}
	offset := fieldv.UnsafeAddr() - valuev.UnsafeAddr()
	if !(0 <= offset && offset < valuev.Type().Size()) {
		panic("provided field is not a member of the given value")
	}

	found := false
	gui(func() {
		tinfo := typeInfo(value)
		for _, engine := range engines {
			fold := engine.values[value]
			for fold != nil {
				found = true
				C.goValueActivate(fold.cvalue, tinfo, C.int(offset))
				fold = fold.next
			}
			// TODO typeNew might also be a linked list keyed by the gvalue.
			//      This would prevent the iteration and the deferrals.
			for fold, _ = range typeNew {
				if fold.gvalue == value {
					found = true
					// Activate these later so they don't get recursively moved
					// out of typeNew while the iteration is still happening.
					defer C.goValueActivate(fold.cvalue, tinfo, C.int(offset))
				}
			}
		}
	})
	if !found {
		// TODO Perhaps return an error instead.
		panic("value is not known")
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
	prev   *valueFold
	next   *valueFold
	owner  valueOwner
}

type valueOwner uint8

const (
	cppOwner = 1 << iota
	jsOwner
)

// wrapGoValue creates a new GoValue object in C++ land wrapping
// the Go value contained in the given interface.
//
// This must be run from the main GUI thread.
func wrapGoValue(engine *Engine, gvalue interface{}, owner valueOwner) (cvalue unsafe.Pointer) {

	// TODO Return an error if gvalue is a non-basic type and not a pointer.
	//      Pointer-to-pointer is also not okay.
	prev, ok := engine.values[gvalue]
	if ok && (prev.owner == owner || owner != cppOwner) {
		return prev.cvalue
	}

	parent := nilPtr
	if owner == cppOwner {
		parent = engine.addr
	}
	fold := &valueFold{
		engine: engine,
		gvalue: gvalue,
		owner:  owner,
	}
	fold.cvalue = C.newGoValue(unsafe.Pointer(fold), typeInfo(gvalue), parent)
	if prev != nil {
		prev.next = fold
		fold.prev = prev
	} else {
		engine.values[gvalue] = fold
	}
	stats.valuesAlive(+1)
	C.engineSetContextForObject(engine.addr, fold.cvalue)
	switch owner {
	case cppOwner:
		C.engineSetOwnershipCPP(engine.addr, fold.cvalue)
	case jsOwner:
		C.engineSetOwnershipJS(engine.addr, fold.cvalue)
	}
	return fold.cvalue
}

// typeNew holds fold values that are created by registered types.
// These values are special in two senses: first, they don't have a
// reference to an engine before they are used in a context that can
// set the reference; second, these values always hold a new cvalue,
// because they are created as a side-effect of the registered type
// being instantiated (it's too late to reuse an existent cvalue).
//
// For these reasons, typeNew holds the fold for these values until
// their engine is known, and once it's known they may have to be
// added to the linked list, since mulitple references for the same
// gvalue may occur.
var typeNew = make(map[*valueFold]bool)

//export hookGoValueTypeNew
func hookGoValueTypeNew(cvalue unsafe.Pointer, specp unsafe.Pointer) (foldp unsafe.Pointer) {
	fold := &valueFold{
		gvalue: (*TypeSpec)(specp).New(),
		cvalue: cvalue,
		owner:  jsOwner,
	}
	typeNew[fold] = true
	stats.valuesAlive(+1)
	return unsafe.Pointer(fold)
}

//export hookGoValueDestroyed
func hookGoValueDestroyed(enginep unsafe.Pointer, foldp unsafe.Pointer) {
	fold := (*valueFold)(foldp)
	engine := fold.engine
	if engine == nil {
		before := len(typeNew)
		delete(typeNew, fold)
		if len(typeNew) == before {
			panic("destroying value without an associated engine; who created the value?")
		}
	} else if engines[engine.addr] == nil {
		// Must never do that. The engine holds memory references that C++ depends on.
		panic(fmt.Sprintf("engine %p was released from global list while its values were still alive", engine.addr))
	} else {
		switch {
		case fold.prev != nil:
			fold.prev.next = fold.next
			if fold.next != nil {
				fold.next.prev = fold.prev
			}
		case fold.next != nil:
			fold.next.prev = fold.prev
			if fold.prev != nil {
				fold.prev.next = fold.next
			} else {
				fold.engine.values[fold.gvalue] = fold.next
			}
		default:
			before := len(engine.values)
			delete(engine.values, fold.gvalue)
			if len(engine.values) == before {
				panic("destroying value that knows about the engine, but the engine doesn't know about the value; who cleared the engine?")
			}
			if engine.destroyed && len(engine.values) == 0 {
				delete(engines, engine.addr)
			}
		}
	}
	stats.valuesAlive(-1)
}

//export hookGoValueReadField
func hookGoValueReadField(enginep, foldp unsafe.Pointer, reflectIndex C.int, resultdv *C.DataValue) {
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
func hookGoValueWriteField(enginep, foldp unsafe.Pointer, reflectIndex C.int, assigndv *C.DataValue) {
	fold := ensureEngine(enginep, foldp)
	v := reflect.ValueOf(fold.gvalue)
	for v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	field := v.Field(int(reflectIndex))
	assign := unpackDataValue(assigndv, fold.engine)

	// TODO Return false to the call site if it fails. That's how Qt seems to handle it internally.
	convertAndSet(field, reflect.ValueOf(assign))
}

func convertAndSet(to, from reflect.Value) {
	defer func() {
		if v := recover(); v != nil {
			// TODO This should be an error. Test and fix.
			panic("FIXME attempted to set a field with the wrong type; this should be an error")
		}
	}()
	to.Set(from.Convert(to.Type()))
}

var (
	dataValueSize  = uintptr(unsafe.Sizeof(C.DataValue{}))
	dataValueArray [C.MaximumParamCount - 1]C.DataValue
)

//export hookGoValueCallMethod
func hookGoValueCallMethod(enginep, foldp unsafe.Pointer, reflectIndex C.int, args *C.DataValue) {
	fold := ensureEngine(enginep, foldp)
	v := reflect.ValueOf(fold.gvalue)

	// TODO Must assert that v is necessarily a pointer here, but we shouldn't have to manipulate
	//      gvalue here for that. This should happen in a sensible place in the wrapping functions
	//      that can still error out to the user in due time.

	method := v.Method(int(reflectIndex))

	// TODO Ensure methods with more parameters than this are not registered.
	var params [C.MaximumParamCount - 1]reflect.Value

	numIn := uintptr(method.Type().NumIn())
	for i := uintptr(0); i < numIn; i++ {
		// TODO Convert the arguments when possible (int32 => int, etc).
		// TODO Type checking to avoid explosions (or catch the explosion)
		paramdv := (*C.DataValue)(unsafe.Pointer(uintptr(unsafe.Pointer(args)) + (i+1)*dataValueSize))
		params[i] = reflect.ValueOf(unpackDataValue(paramdv, fold.engine))
	}

	result := method.Call(params[:numIn])

	if len(result) == 1 {
		packDataValue(result[0].Interface(), args, fold.engine, jsOwner)
	} else if len(result) > 1 {
		if len(result) > len(dataValueArray) {
			panic("function has too many results")
		}
		for i, v := range result {
			packDataValue(v.Interface(), &dataValueArray[i], fold.engine, jsOwner)
		}
		args.dataType = C.DTList
		*(*unsafe.Pointer)(unsafe.Pointer(&args.data)) = C.newVariantList(&dataValueArray[0], C.int(len(result)))
	}
}

func ensureEngine(enginep, foldp unsafe.Pointer) *valueFold {
	fold := (*valueFold)(foldp)
	if fold.engine != nil {
		return fold
	}

	if enginep == nilPtr {
		panic("accessing value without an engine pointer; who created the value?")
	}
	engine := engines[enginep]
	if engine == nil {
		panic("unknown engine pointer; who created the engine?")
	}
	fold.engine = engine
	prev := engine.values[fold.gvalue]
	if prev != nil {
		for prev.next != nil {
			prev = prev.next
		}
		prev.next = fold
		fold.prev = prev
	} else {
		engine.values[fold.gvalue] = fold
	}
	before := len(typeNew)
	delete(typeNew, fold)
	if len(typeNew) == before {
		panic("value had no engine, but was not created by a registered type; who created the value?")
	}
	return fold
}

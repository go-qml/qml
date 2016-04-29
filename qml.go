package qml

// #include <stdlib.h>
//
// #include "capi.h"
//
import "C"

import (
	"errors"
	"unsafe"
)

func AddLibraryPath(path string) {
	cpath, cpathLen := unsafeStringData(path)
	RunMain(func() {
		C.coreAddLibraryPath(cpath, cpathLen)
	})
}

func cerror(cerr *C.error) error {
	err := errors.New(C.GoString((*C.char)(unsafe.Pointer(cerr))))
	C.free(unsafe.Pointer(cerr))
	return err
}

func cmust(cerr *C.error) {
	if cerr != nil {
		panic(cerror(cerr).Error())
	}
}

// TODO Signal emitting support for go values.

// LoadResources registers all resources in the provided resources collection,
// making them available to be loaded by any Engine and QML file.
// Registered resources are made available under "qrc:///some/path", where
// "some/path" is the path the resource was added with.
func LoadResources(r *Resources) {
	var base unsafe.Pointer
	if len(r.sdata) > 0 {
		base = *(*unsafe.Pointer)(unsafe.Pointer(&r.sdata))
	} else if len(r.bdata) > 0 {
		base = *(*unsafe.Pointer)(unsafe.Pointer(&r.bdata))
	}
	tree := (*C.char)(unsafe.Pointer(uintptr(base) + uintptr(r.treeOffset)))
	name := (*C.char)(unsafe.Pointer(uintptr(base) + uintptr(r.nameOffset)))
	data := (*C.char)(unsafe.Pointer(uintptr(base) + uintptr(r.dataOffset)))
	C.registerResourceData(C.int(r.version), tree, name, data)
}

// UnloadResources unregisters all previously registered resources from r.
func UnloadResources(r *Resources) {
	var base unsafe.Pointer
	if len(r.sdata) > 0 {
		base = *(*unsafe.Pointer)(unsafe.Pointer(&r.sdata))
	} else if len(r.bdata) > 0 {
		base = *(*unsafe.Pointer)(unsafe.Pointer(&r.bdata))
	}
	tree := (*C.char)(unsafe.Pointer(uintptr(base) + uintptr(r.treeOffset)))
	name := (*C.char)(unsafe.Pointer(uintptr(base) + uintptr(r.nameOffset)))
	data := (*C.char)(unsafe.Pointer(uintptr(base) + uintptr(r.dataOffset)))
	C.unregisterResourceData(C.int(r.version), tree, name, data)
}

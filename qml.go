package qml

// #include <stdlib.h>
//
// #include "capi.h"
//
import "C"

import (
	"errors"
	"unsafe"

	"github.com/limetext/qml-go/internal/util"
)

func AddLibraryPath(path string) {
	cpath, cpathLen := util.UnsafeStringData(path)
	RunMain(func() {
		C.coreAddLibraryPath((*C.char)(cpath), C.int(cpathLen))
	})
}

func SetApplicationName(name string) {
	cname, cnameLen := util.UnsafeStringData(name)
	RunMain(func() {
		C.coreSetApplicationName((*C.char)(cname), C.int(cnameLen))
	})
}

func SetApplicationVersion(version string) {
	cversion, cversionLen := util.UnsafeStringData(version)
	RunMain(func() {
		C.coreSetApplicationVersion((*C.char)(cversion), C.int(cversionLen))
	})
}

func SetOrganizationDomain(domain string) {
	cdomain, cdomainLen := util.UnsafeStringData(domain)
	RunMain(func() {
		C.coreSetOrganizationDomain((*C.char)(cdomain), C.int(cdomainLen))
	})
}

func SetOrganizationName(name string) {
	cname, cnameLen := util.UnsafeStringData(name)
	RunMain(func() {
		C.coreSetOrganizationName((*C.char)(cname), C.int(cnameLen))
	})
}

func SetApplicationDisplayName(name string) {
	cname, cnameLen := util.UnsafeStringData(name)
	RunMain(func() {
		C.guiappSetApplicationDisplayName((*C.char)(cname), C.int(cnameLen))
	})
}

// Qt 5.7
// func SetDesktopFileName(path string) {
// 	cpath, cpathLen := util.UnsafeStringData(path)
// 	RunMain(func() {
// 		C.guiappSetDesktopFileName((*C.char)(cpath), C.int(cpathLen))
// 	})
// }

func SetWindowIcon(path string) {
	cpath, cpathLen := util.UnsafeStringData(path)
	RunMain(func() {
		C.guiappSetWindowIcon((*C.char)(cpath), C.int(cpathLen))
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

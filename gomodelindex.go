package qml

// #cgo CPPFLAGS: -I../cpp
// #cgo CXXFLAGS: -std=c++0x -pedantic-errors -Wall -fno-strict-aliasing
// #cgo LDFLAGS: -lstdc++
// #cgo pkg-config: Qt5Core Qt5Widgets Qt5Quick
//
// #include <stdlib.h>
// #include "goitemmodel_api.h"
//
import "C"
import "unsafe"

type qModelIndex struct {
	ptr    uintptr
	engine *Engine
}

func mkModelIndex(ptr uintptr, engine *Engine) ModelIndex {
	if ptr == 0 {
		return nil
	}
	return &qModelIndex{
		ptr:    ptr,
		engine: engine,
	}
}

type ModelIndex interface {
	// ModelIndex can only be created from a ItemModel
	internal_ModelIndex()
	Child(row, col int) ModelIndex
	Sibling(row, col int) ModelIndex
	Column() int
	Row() int
	Data(role Role) interface{}
	Flags() ItemFlags
	InternalId() uintptr
	InternalPointer() uintptr
	IsValid() bool
	Model() ItemModel
	Parent() ModelIndex
}

func (i *qModelIndex) internal_ModelIndex() {}

func (i *qModelIndex) Child(row, col int) ModelIndex {
	return mkModelIndex(uintptr(C.modelIndexChild(unsafe.Pointer(i.ptr), C.int(row), C.int(col))), i.engine)
}

func (i *qModelIndex) Sibling(row, col int) ModelIndex {
	return mkModelIndex(uintptr(C.modelIndexSibling(unsafe.Pointer(i.ptr), C.int(row), C.int(col))), i.engine)
}

func (i *qModelIndex) Column() int {
	return int(C.modelIndexColumn(unsafe.Pointer(i.ptr)))
}

func (i *qModelIndex) Row() int {
	return int(C.modelIndexRow(unsafe.Pointer(i.ptr)))
}

func (i *qModelIndex) Data(role Role) interface{} {
	var dvalue C.DataValue
	C.modelIndexData(unsafe.Pointer(i.ptr), C.int(role), &dvalue)
	return unpackDataValue(&dvalue, i.engine)
}

func (i *qModelIndex) Flags() ItemFlags {
	return ItemFlags(C.modelIndexFlags(unsafe.Pointer(i.ptr)))
}

func (i *qModelIndex) InternalId() uintptr {
	return uintptr(C.modelIndexInternalId(unsafe.Pointer(i.ptr)))
}

func (i *qModelIndex) InternalPointer() uintptr {
	return uintptr(C.modelIndexInternalPointer(unsafe.Pointer(i.ptr)))
}

func (i *qModelIndex) IsValid() bool {
	if i == nil || i.ptr == 0 {
		return false
	}
	return bool(C.modelIndexIsValid(unsafe.Pointer(i.ptr)))
}

func (i *qModelIndex) Model() ItemModel {
	return itemModelFromCPP(uintptr(C.modelIndexModel(unsafe.Pointer(i.ptr))), i.engine)
}

func (i *qModelIndex) Parent() ModelIndex {
	return mkModelIndex(uintptr(C.modelIndexParent(unsafe.Pointer(i.ptr))), i.engine)
}

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
import (
	"runtime"
	"unsafe"
)

const (
	NoItemFlags          ItemFlags = 0
	ItemIsSelectable     ItemFlags = 1
	ItemIsEditable       ItemFlags = 2
	ItemIsDragEnabled    ItemFlags = 4
	ItemIsDropEnabled    ItemFlags = 8
	ItemIsUserCheckable  ItemFlags = 16
	ItemIsEnabled        ItemFlags = 32
	ItemIsAutoTristate   ItemFlags = 64
	ItemNeverHasChildren ItemFlags = 128
	ItemIsUserTristate   ItemFlags = 256
)

type ItemModel interface {
	internal_ItemModel() *goItemModel
}

type ItemModelInternal interface {
	CreateIndex(row int, column int, id uintptr) ModelIndex

	BeginInsertColumns(parent ModelIndex, first int, last int)
	EndInsertColumns()
	BeginInsertRows(parent ModelIndex, first int, last int)
	EndInsertRows()
	BeginRemoveRows(parent ModelIndex, first int, last int)
	EndRemoveRows()
	DataChanged(topLeft ModelIndex, bottomRight ModelIndex)
}

type goItemModel struct {
	common *Common
	impl   ItemModelImpl
}

func (q *goItemModel) internal_ItemModel() *goItemModel { return q }

func NewItemModel(engine *Engine, parent Object, impl ItemModelImpl) (ItemModel, ItemModelInternal) {
	im := mkItemModel()

	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = unsafe.Pointer(parent.Addr())
	}

	im.impl = impl

	var imptr unsafe.Pointer

	RunMain(func() {
		imptr = C.newGoItemModel(parentPtr, unsafe.Pointer(im))
	})

	im.common = CommonOf(imptr, engine)

	return im, im
}

func mkItemModel() *goItemModel {
	qimg := &goItemModel{}

	runtime.SetFinalizer(qimg, deleteItemModel)

	return qimg
}

func deleteItemModel(qim *goItemModel) {
	C.deleteGoItemModel(qim.common.addr)
	// qim.addr = nil
}

type ItemModelImpl interface {
	// Required functions
	ColumnCount(parent ModelIndex) int
	Data(index ModelIndex, role Role) interface{}
	Index(row int, column int, parent ModelIndex) ModelIndex
	Parent(index ModelIndex) ModelIndex
	RowCount(parent ModelIndex) int

	// Required for editing
	Flags(index ModelIndex) ItemFlags
	SetData(index ModelIndex, value interface{}, role Role) bool
}

type Role int
type ItemFlags int

type cppItemModelImpl struct {
	qim *goItemModel
}

func itemModelFromCPP(addr uintptr, engine *Engine) ItemModel {
	im := mkItemModel()

	im.impl = &cppItemModelImpl{im}
	im.common = CommonOf(unsafe.Pointer(addr), engine)

	return im
}

func passMI(mi ModelIndex) unsafe.Pointer {
	if mi == nil {
		return nil
	}
	return unsafe.Pointer(mi.(*qModelIndex).ptr)
}

// Required functions
func (cim *cppItemModelImpl) ColumnCount(parent ModelIndex) int {
	return int(C.itemModelColumnCount(cim.qim.common.addr, passMI(parent)))
}

func (cim *cppItemModelImpl) Data(index ModelIndex, role Role) interface{} {
	var dvalue C.DataValue
	C.itemModelData(cim.qim.common.addr, passMI(index), C.int(role), &dvalue)
	return unpackDataValue(&dvalue, cim.qim.common.engine)
}

func (cim *cppItemModelImpl) Index(row int, col int, parent ModelIndex) ModelIndex {
	return mkModelIndex(uintptr(C.itemModelIndex(cim.qim.common.addr, C.int(row), C.int(col), passMI(parent))), cim.qim.common.engine)
}

func (cim *cppItemModelImpl) Parent(index ModelIndex) ModelIndex {
	return mkModelIndex(uintptr(C.itemModelParent(cim.qim.common.addr, passMI(index))), cim.qim.common.engine)
}

func (cim *cppItemModelImpl) RowCount(parent ModelIndex) int {
	return int(C.itemModelRowCount(cim.qim.common.addr, passMI(parent)))
}

// Required for editing
func (cim *cppItemModelImpl) Flags(index ModelIndex) ItemFlags {
	return ItemFlags(C.itemModelFlags(cim.qim.common.addr, passMI(index)))
}

func (cim *cppItemModelImpl) SetData(index ModelIndex, value interface{}, role Role) bool {
	var dvalue C.DataValue
	packDataValue(value, &dvalue, cim.qim.common.engine, cppOwner)
	return bool(C.itemModelSetData(cim.qim.common.addr, passMI(index), &dvalue, C.int(role)))
}

// Internal (protected) functions

func (qim *goItemModel) CreateIndex(row int, column int, id uintptr) ModelIndex {
	indexPtr := C.itemModelCreateIndex(qim.common.addr, C.int(row), C.int(column), C.uintptr_t(id))
	return mkModelIndex(uintptr(indexPtr), qim.common.engine)
}

func (qim *goItemModel) BeginInsertColumns(parent ModelIndex, first int, last int) {
	C.itemModelBeginInsertColumns(qim.common.addr, passMI(parent), C.int(first), C.int(last))
}

func (qim *goItemModel) EndInsertColumns() {
	C.itemModelEndInsertColumns(qim.common.addr)
}

func (qim *goItemModel) BeginInsertRows(parent ModelIndex, first int, last int) {
	C.itemModelBeginInsertRows(qim.common.addr, passMI(parent), C.int(first), C.int(last))
}

func (qim *goItemModel) EndInsertRows() {
	C.itemModelEndInsertRows(qim.common.addr)
}

func (qim *goItemModel) BeginRemoveRows(parent ModelIndex, first int, last int) {
	C.itemModelBeginRemoveRows(qim.common.addr, passMI(parent), C.int(first), C.int(last))
}

func (qim *goItemModel) EndRemoveRows() {
	C.itemModelEndRemoveRows(qim.common.addr)
}

func (qim *goItemModel) DataChanged(topLeft ModelIndex, bottomRight ModelIndex) {
	C.itemModelDataChanged(qim.common.addr, passMI(topLeft), passMI(bottomRight))
}

// Required functions
//export implColumnCount
func implColumnCount(qim uintptr, parent uintptr) int {
	im := (*goItemModel)(unsafe.Pointer(qim))
	parentMi := mkModelIndex(parent, im.common.engine)

	return im.impl.ColumnCount(parentMi)
}

//export implData
func implData(qim uintptr, index uintptr, role int, ret *C.DataValue) {
	im := (*goItemModel)(unsafe.Pointer(qim))
	indexMi := mkModelIndex(index, im.common.engine)

	v := im.impl.Data(indexMi, Role(role))

	packDataValue(v, ret, im.common.engine, cppOwner)
}

//export implIndex
func implIndex(qim uintptr, row int, column int, parent uintptr) uintptr {
	im := (*goItemModel)(unsafe.Pointer(qim))
	parentMi := mkModelIndex(parent, im.common.engine)

	ret := im.impl.Index(row, column, parentMi)
	if ret != nil {
		return ret.(*qModelIndex).ptr
	}

	return 0
}

//export implParent
func implParent(qim uintptr, index uintptr) uintptr {
	im := (*goItemModel)(unsafe.Pointer(qim))
	indexMi := mkModelIndex(index, im.common.engine)
	parentMi := im.impl.Parent(indexMi)
	if parentMi != nil {
		return parentMi.(*qModelIndex).ptr
	}

	return 0
}

//export implRowCount
func implRowCount(qim uintptr, parent uintptr) int {
	im := (*goItemModel)(unsafe.Pointer(qim))
	parentMi := mkModelIndex(parent, im.common.engine)

	return im.impl.RowCount(parentMi)
}

// Required for editing
//export implFlags
func implFlags(qim uintptr, index uintptr) ItemFlags {
	im := (*goItemModel)(unsafe.Pointer(qim))
	indexMi := mkModelIndex(index, im.common.engine)

	return im.impl.Flags(indexMi)
}

//export implSetData
func implSetData(qim uintptr, index uintptr, dv *C.DataValue, role int) bool {
	im := (*goItemModel)(unsafe.Pointer(qim))
	indexMi := mkModelIndex(index, im.common.engine)

	value := unpackDataValue(dv, im.common.engine)

	return im.impl.SetData(indexMi, value, Role(role))
}

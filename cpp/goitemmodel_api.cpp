
#include "goitemmodel.h"
#include "goitemmodel_api.h"
#include "util.cpp"
#include <QAbstractItemModel>


QItemModel_* newGoItemModel(QObject_* parent, GoAddr* impl) {
	return reinterpret_cast<QItemModel_*>(new GoItemModel(reinterpret_cast<QObject*>(parent), impl));
}

void deleteGoItemModel(QItemModel_* im) {
	delete reinterpret_cast<GoItemModel*>(im);
}



QModelIndex_ *modelIndexChild(QModelIndex_ *mi, int row, int col) {
	QModelIndex *ret = new QModelIndex;
	*ret = reinterpret_cast<QModelIndex*>(mi)->child(row, col);
	return (QModelIndex_*)ret;
}

QModelIndex_ *modelIndexSibling(QModelIndex_ *mi, int row, int col) {
	QModelIndex *ret = new QModelIndex;
	*ret = reinterpret_cast<QModelIndex*>(mi)->sibling(row, col);
	return (QModelIndex_*)ret;
}

int modelIndexColumn(QModelIndex_ *mi) {
	return reinterpret_cast<QModelIndex*>(mi)->column();
}

int modelIndexRow(QModelIndex_ *mi) {
	return reinterpret_cast<QModelIndex*>(mi)->row();
}

void modelIndexData(QModelIndex_ *mi, int role, DataValue *ret) {
	QVariant value = reinterpret_cast<QModelIndex*>(mi)->data(role);
  packDataValue(&value, ret);
}

quint32 modelIndexFlags(QModelIndex_ *mi) {
	return reinterpret_cast<QModelIndex*>(mi)->flags();
}

uintptr_t modelIndexInternalId(QModelIndex_ *mi) {
	return reinterpret_cast<QModelIndex*>(mi)->internalId();
}

uintptr_t modelIndexInternalPointer(QModelIndex_ *mi) {
	return (quintptr)reinterpret_cast<QModelIndex*>(mi)->internalPointer();
}

bool modelIndexIsValid(QModelIndex_ *mi) {
	return reinterpret_cast<QModelIndex*>(mi)->isValid();
}

QItemModel_ *modelIndexModel(QModelIndex_ *mi) {
	return (QItemModel_ *)reinterpret_cast<QModelIndex*>(mi)->model();
}

QModelIndex_ *modelIndexParent(QModelIndex_ *mi) {
	QModelIndex *ret = new QModelIndex;
	*ret = reinterpret_cast<QModelIndex*>(mi)->parent();
	return (QModelIndex_*)ret;
}


// Protected functions

QModelIndex_ *itemModelCreateIndex(QItemModel_ *im, int row, int col, uintptr_t id) {
	QModelIndex *ret = new QModelIndex;
	*ret = reinterpret_cast<GoItemModel*>(im)->protCreateIndex(row, col, id);
	return (QModelIndex_*)ret;
}

void itemModelBeginInsertColumns(QItemModel_ *im, QModelIndex_ *parent, int first, int last) {
	reinterpret_cast<GoItemModel*>(im)->protBeginInsertColumns(miCastFrom(parent), first, last);
}

void itemModelEndInsertColumns(QItemModel_ *im) {
	reinterpret_cast<GoItemModel*>(im)->protEndInsertColumns();
}

void itemModelBeginInsertRows(QItemModel_ *im, QModelIndex_ *parent, int first, int last) {
	reinterpret_cast<GoItemModel*>(im)->protBeginInsertRows(miCastFrom(parent), first, last);
}

void itemModelEndInsertRows(QItemModel_ *im) {
	reinterpret_cast<GoItemModel*>(im)->protEndInsertRows();
}

void itemModelDataChanged(QItemModel_ *im, QModelIndex_ *topLeft, QModelIndex_ *bottomRight) {
	reinterpret_cast<GoItemModel*>(im)->dataChanged(miCastFrom(topLeft), miCastFrom(bottomRight));
}

// Required functions
int itemModelColumnCount(QItemModel_ *im, QModelIndex_ *parent) {
	return reinterpret_cast<QAbstractItemModel*>(im)->columnCount(miCastFrom(parent));
}

void itemModelData(QItemModel_ *im, QModelIndex_ *index, int role, DataValue *ret) {
	QVariant var = reinterpret_cast<QAbstractItemModel*>(im)->data(miCastFrom(index), role);

  packDataValue(&var, ret);
}

QModelIndex_ *itemModelIndex(QItemModel_ *im, int row, int column, QModelIndex_ *parent) {
	QModelIndex *ret = new QModelIndex;
	*ret = reinterpret_cast<QAbstractItemModel*>(im)->index(row, column, miCastFrom(parent));
	return (QModelIndex_*)ret;
}

QModelIndex_ *itemModelParent(QItemModel_ *im, QModelIndex_ *index) {
	QModelIndex *ret = new QModelIndex;
	*ret = reinterpret_cast<QAbstractItemModel*>(im)->parent(miCastFrom(index));
	return (QModelIndex_*)ret;
}

int itemModelRowCount(QItemModel_ *im, QModelIndex_ *parent) {
	return reinterpret_cast<QAbstractItemModel*>(im)->rowCount(miCastFrom(parent));
}


// Required for editing
int itemModelFlags(QItemModel_ *im, QModelIndex_ *index) {
	return reinterpret_cast<QAbstractItemModel*>(im)->flags(miCastFrom(index));
}

bool itemModelSetData(QItemModel_ *im, QModelIndex_ *index, DataValue *value, int role) {
	QVariant var;
  unpackDataValue(value, &var);
	return reinterpret_cast<QAbstractItemModel*>(im)->setData(miCastFrom(index), var, role);
}

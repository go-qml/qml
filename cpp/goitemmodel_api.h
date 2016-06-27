#ifndef GOITEMMODEL_API_H
#define GOITEMMODEL_API_H

// #include "capi.h"
#include <stdint.h>
#include <stddef.h>
#include <stdbool.h>
#include "capi.h"
#include <QtGlobal>

typedef void GoAddr;
typedef void QObject_;
typedef void QItemModel_;
typedef void QModelIndex_;
typedef void QVariant_;

#ifdef __cplusplus
extern "C" {
#endif

QItemModel_* newGoItemModel(QObject_*, GoAddr*);
void deleteGoItemModel(QItemModel_*);

QModelIndex_ *modelIndexChild(QModelIndex_ *mi, int row, int col);
QModelIndex_ *modelIndexSibling(QModelIndex_ *mi, int row, int col);
int modelIndexColumn(QModelIndex_ *mi);
int modelIndexRow(QModelIndex_ *mi);
void modelIndexData(QModelIndex_ *mi, int role, DataValue *ret);
quint32 modelIndexFlags(QModelIndex_ *mi);
uintptr_t modelIndexInternalId(QModelIndex_ *mi);
uintptr_t modelIndexInternalPointer(QModelIndex_ *mi);
bool modelIndexIsValid(QModelIndex_ *mi);
QItemModel_ *modelIndexModel(QModelIndex_ *mi);
QModelIndex_ *modelIndexParent(QModelIndex_ *mi);


// Required functions
int itemModelColumnCount(QItemModel_ *im, QModelIndex_ *parent);
void itemModelData(QItemModel_ *im, QModelIndex_ *index, int role, DataValue *ret);
QModelIndex_ *itemModelIndex(QItemModel_ *im, int row, int column, QModelIndex_ *parent);
QModelIndex_ *itemModelParent(QItemModel_ *im, QModelIndex_ *index);
int itemModelRowCount(QItemModel_ *im, QModelIndex_ *parent);

// Required for editing
int itemModelFlags(QItemModel_ *im, QModelIndex_ *index);
bool itemModelSetData(QItemModel_ *im, QModelIndex_ *index, DataValue *value, int role);

// Protected functions
QModelIndex_ *itemModelCreateIndex(QItemModel_ *im, int row, int col, uintptr_t id);
void itemModelBeginInsertColumns(QItemModel_ *im, QModelIndex_ *parent, int first, int last);
void itemModelEndInsertColumns(QItemModel_ *im);
void itemModelBeginInsertRows(QItemModel_ *im, QModelIndex_ *parent, int first, int last);
void itemModelEndInsertRows(QItemModel_ *im);


void itemModelDataChanged(QItemModel_ *im, QModelIndex_ *topLeft, QModelIndex_ *bottomRight);


// virtual QModelIndex 	buddy(const QModelIndex &index) const
// virtual bool 	canDropMimeData(const QMimeData *data, Qt::DropAction action, int row, int column, const QModelIndex &parent) const
// virtual bool 	canFetchMore(const QModelIndex &parent) const
// virtual bool 	dropMimeData(const QMimeData *data, Qt::DropAction action, int row, int column, const QModelIndex &parent)
// virtual void 	fetchMore(const QModelIndex &parent)
// virtual Qt::ItemFlags 	flags(const QModelIndex &index) const
// virtual bool 	hasChildren(const QModelIndex &parent = QModelIndex()) const
// bool 	hasIndex(int row, int column, const QModelIndex &parent = QModelIndex()) const
// virtual QVariant 	headerData(int section, Qt::Orientation orientation, int role = Qt::DisplayRole) const
// bool 	insertColumn(int column, const QModelIndex &parent = QModelIndex())
// virtual bool 	insertColumns(int column, int count, const QModelIndex &parent = QModelIndex())
// bool 	insertRow(int row, const QModelIndex &parent = QModelIndex())
// virtual bool 	insertRows(int row, int count, const QModelIndex &parent = QModelIndex())
// virtual QMap<int, QVariant> 	itemData(const QModelIndex &index) const
// virtual QModelIndexList 	match(const QModelIndex &start, int role, const QVariant &value, int hits = 1, Qt::MatchFlags flags = Qt::MatchFlags( Qt::MatchStartsWith | Qt::MatchWrap )) const
// virtual QMimeData *	mimeData(const QModelIndexList &indexes) const
// virtual QStringList 	mimeTypes() const
// bool 	moveColumn(const QModelIndex &sourceParent, int sourceColumn, const QModelIndex &destinationParent, int destinationChild)
// virtual bool 	moveColumns(const QModelIndex &sourceParent, int sourceColumn, int count, const QModelIndex &destinationParent, int destinationChild)
// bool 	moveRow(const QModelIndex &sourceParent, int sourceRow, const QModelIndex &destinationParent, int destinationChild)
// virtual bool 	moveRows(const QModelIndex &sourceParent, int sourceRow, int count, const QModelIndex &destinationParent, int destinationChild)
// bool 	removeColumn(int column, const QModelIndex &parent = QModelIndex())
// virtual bool 	removeColumns(int column, int count, const QModelIndex &parent = QModelIndex())
// bool 	removeRow(int row, const QModelIndex &parent = QModelIndex())
// virtual bool 	removeRows(int row, int count, const QModelIndex &parent = QModelIndex())
// virtual QHash<int, QByteArray> 	roleNames() const
// virtual bool 	setData(const QModelIndex &index, const QVariant &value, int role = Qt::EditRole)
// virtual bool 	setHeaderData(int section, Qt::Orientation orientation, const QVariant &value, int role = Qt::EditRole)
// virtual bool 	setItemData(const QModelIndex &index, const QMap<int, QVariant> &roles)
// virtual QModelIndex 	sibling(int row, int column, const QModelIndex &index) const
// virtual void 	sort(int column, Qt::SortOrder order = Qt::AscendingOrder)
// virtual QSize 	span(const QModelIndex &index) const
// virtual Qt::DropActions 	supportedDragActions() const
// virtual Qt::DropActions 	supportedDropActions() const


#ifdef __cplusplus
} // extern "C"
#endif

#endif // GOITEMMODEL_API_H

// vim:ts=4:sw=4:et:ft=cpp

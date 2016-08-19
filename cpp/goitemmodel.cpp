
#include "goitemmodel.h"

GoItemModel::GoItemModel(QObject* parent, GoValueRef valueref)
	: QAbstractItemModel(parent), valueref(valueref) {

}

QModelIndex_* miCastTo(const QModelIndex &index) {
	return reinterpret_cast<QModelIndex_*>(const_cast<QModelIndex*>(&index));
}
QModelIndex miCastFrom(QModelIndex_ *index) {
	if (index == NULL) {
		return QModelIndex();
	}
	return *reinterpret_cast<QModelIndex*>(index);
}

// Required functions
int GoItemModel::columnCount(const QModelIndex &parent) const {
	return implColumnCount(valueref, miCastTo(parent));
}

QVariant GoItemModel::data(const QModelIndex &index, int role) const {
	DataValue value;

	implData(valueref, miCastTo(index), role, &value);

  QVariant var;
  unpackDataValue(&value, &var);

	return var;
}

QModelIndex GoItemModel::index(int row, int column, const QModelIndex &parent) const {
	return miCastFrom(implIndex(valueref, row, column, miCastTo(parent)));
}

QModelIndex GoItemModel::parent(const QModelIndex &index) const {
	return miCastFrom(implParent(valueref, miCastTo(index)));
}

int GoItemModel::rowCount(const QModelIndex &parent) const {
	return implRowCount(valueref, miCastTo(parent));
}


// Required for editing
Qt::ItemFlags GoItemModel::flags(const QModelIndex &index) const {
	return (Qt::ItemFlags)implFlags(valueref, miCastTo(index));
}

bool GoItemModel::setData(const QModelIndex &index, const QVariant &value, int role) {
	DataValue *dv = (DataValue *) malloc(sizeof(DataValue));
    packDataValue(&value, dv);
	return implSetData(valueref, miCastTo(index), dv, role);
}

// Internal Protected functions

QModelIndex GoItemModel::protCreateIndex(int row, int column, quintptr id) const {
	return createIndex(row, column, id);
}

void GoItemModel::protBeginInsertColumns(const QModelIndex &parent, int first, int last) {
	return beginInsertColumns(parent, first, last);
}

void GoItemModel::protEndInsertColumns() {
	return endInsertColumns();
}

void GoItemModel::protBeginInsertRows(const QModelIndex &parent, int first, int last) {
	return beginInsertRows(parent, first, last);
}

void GoItemModel::protEndInsertRows() {
	return endInsertRows();
}

void GoItemModel::protBeginRemoveRows(const QModelIndex &parent, int first, int last) {
	return beginRemoveRows(parent, first, last);
}

void GoItemModel::protEndRemoveRows() {
	return endRemoveRows();
}

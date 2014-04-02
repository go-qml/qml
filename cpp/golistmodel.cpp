#include "golistmodel.h"

#include <QtCore/qvector.h>

#include <algorithm>

GoListModel::GoListModel(const QVariantList &list, GoTypeInfo *typeInfo, QObject *parent)
        : QAbstractListModel(parent), typeInfo(typeInfo), lst(list)
{
    if (typeInfo) {
        GoMemberInfo *memberInfo = typeInfo->fields;
        for (int i = 0; i < typeInfo->fieldsLen; i++) {
            roles[Qt::UserRole + i] = memberInfo->memberName;
            memberInfo++;
        }
    }
}

int GoListModel::rowCount(const QModelIndex &parent) const
{
    if (parent.isValid()) {
        return 0;
    }
    return lst.count();
}

QModelIndex GoListModel::sibling(int row, int column, const QModelIndex &idx) const
{
    if (!idx.isValid() || column != 0 || row >= lst.count()) {
        return QModelIndex();
    }
    return createIndex(row, 0);
}

QVariant GoListModel::data(const QModelIndex &index, int role) const
{
    if (index.row() < 0 || index.row() >= lst.size()) {
        return QVariant();
    }
    QVariant var = lst.at(index.row());
    if (role == Qt::DisplayRole || role == Qt::EditRole) {
        return var;
    }
    if (role >= Qt::UserRole && typeInfo != NULL && role-Qt::UserRole < typeInfo->fieldsLen) {
        GoMemberInfo *memberInfo = typeInfo->fields + (role-Qt::UserRole);
        return var.value<QObject*>()->property(memberInfo->memberName);
    }
    return QVariant();
}

Qt::ItemFlags GoListModel::flags(const QModelIndex &index) const
{
    if (!index.isValid()) {
        return QAbstractListModel::flags(index) | Qt::ItemIsDropEnabled;
    }
    return QAbstractListModel::flags(index) | Qt::ItemIsEditable | Qt::ItemIsDragEnabled | Qt::ItemIsDropEnabled;
}

bool GoListModel::setData(const QModelIndex &index, const QVariant &value, int role)
{
    if (index.row() >= 0 && index.row() < lst.size() && (role == Qt::EditRole || role == Qt::DisplayRole)) {
        lst.replace(index.row(), value);
        emit dataChanged(index, index, QVector<int>() << role);
        return true;
    }
    return false;
}

bool GoListModel::insertRows(int row, int count, const QModelIndex &parent)
{
    if (count < 1 || row < 0 || row > rowCount(parent))
        return false;

    beginInsertRows(QModelIndex(), row, row + count - 1);
    for (int r = 0; r < count; ++r)
        lst.insert(row, QVariant());
    endInsertRows();

    return true;
}

bool GoListModel::removeRows(int row, int count, const QModelIndex &parent)
{
    if (count <= 0 || row < 0 || (row + count) > rowCount(parent))
        return false;

    beginRemoveRows(QModelIndex(), row, row + count - 1);

    for (int r = 0; r < count; ++r)
        lst.removeAt(row);

    endRemoveRows();
    return true;
}

QVariantList GoListModel::list() const
{
    return lst;
}

void GoListModel::setList(const QVariantList &list)
{
    emit beginResetModel();
    lst = list;
    emit endResetModel();
}

QHash<int,QByteArray> GoListModel::roleNames() const
{
	return roles;
}

// vim:ts=4:sw=4:et:ft=cpp

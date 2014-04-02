#ifndef GOLISTMODEL_H
#define GOLISTMODEL_H

#include <QAbstractListModel>


class Q_CORE_EXPORT GoListModel : public QAbstractListModel
{
    Q_OBJECT

public:

    explicit GoListModel(QObject *parent = 0) : QAbstractListModel(parent) {};
    explicit GoListModel(const QVariantList &list, GoTypeInfo *typeInfo, QObject *parent = 0);

    int rowCount(const QModelIndex &parent = QModelIndex()) const;
    QModelIndex sibling(int row, int column, const QModelIndex &idx) const;

    QVariant data(const QModelIndex &index, int role) const;
    bool setData(const QModelIndex &index, const QVariant &value, int role = Qt::EditRole);

    Qt::ItemFlags flags(const QModelIndex &index) const;

    bool insertRows(int row, int count, const QModelIndex &parent = QModelIndex());
    bool removeRows(int row, int count, const QModelIndex &parent = QModelIndex());

    virtual QHash<int,QByteArray> roleNames() const;

    //void sort(int column, Qt::SortOrder order = Qt::AscendingOrder);

    QVariantList list() const;
    void setList(const QVariantList &list);

    //Qt::DropActions supportedDropActions() const;

private:

    Q_DISABLE_COPY(GoListModel)

    GoTypeInfo *typeInfo;

    QHash<int,QByteArray> roles;

    QVariantList lst;
};

#endif // GOLISTMODEL_H

// vim:ts=4:sw=4:et:ft=cpp

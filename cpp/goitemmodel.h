#ifndef GOITEMMODEL_H
#define GOITEMMODEL_H

// #include "capi.h"
#include <QObject>
#include <QAbstractItemModel>
#include "goitemmodel_impl.h"

// QImage_* newQImage(int width, int height, unsigned int format);
// QImage_* loadQImage(const char *filename, int filename_length, const char *format);
// void deleteQImage(QImage_*);

class GoItemModel : public QAbstractItemModel
{
  Q_OBJECT
public:
 	// QAbstractItemModel(QObject *parent = Q_NULLPTR)
 	GoItemModel(QObject *parent, GoAddr *impl);
	// virtual 	~QAbstractItemModel()

	// Required functions
	int 			columnCount(const QModelIndex &parent = QModelIndex()) const;
	QVariant 		data(const QModelIndex &index, int role = Qt::DisplayRole) const override;
	QModelIndex 	index(int row, int column, const QModelIndex &parent = QModelIndex()) const;
	QModelIndex 	parent(const QModelIndex &index) const;
	int 			rowCount(const QModelIndex &parent = QModelIndex()) const;

	// Required for editing
	Qt::ItemFlags 	flags(const QModelIndex &index) const;
	bool 			setData(const QModelIndex &index, const QVariant &value, int role = Qt::EditRole);

  QModelIndex 	protCreateIndex(int row, int column, quintptr id) const;
  void protBeginInsertColumns(const QModelIndex &parent, int first, int last);
  void protEndInsertColumns();
  void protBeginInsertRows(const QModelIndex &parent, int first, int last);
  void protEndInsertRows();



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

// private:
    GoAddr *addr;
};

#endif // GOITEMMODEL_H

// vim:ts=4:sw=4:et:ft=cpp

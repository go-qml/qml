
#include "qmodelindex.h"
#include "util.cpp"


QModelIndex_* newGoModelIndex(QObject_* parent, GoValueRef impl) {
	return reinterpret_cast<QModelIndex_*>(new GoModelIndex(reinterpret_cast<QObject*>(parent), impl));
}

void deleteGoModelIndex(QModelIndex_* im) {
	delete reinterpret_cast<QModelIndex*>(im);
}


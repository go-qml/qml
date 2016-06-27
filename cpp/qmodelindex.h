#ifndef MODELINDEX_H
#define MODELINDEX_H

// #include "capi.h"

typedef void GoAddr;
typedef void QObject_;
typedef void QItemModel_;
typedef void QModelIndex_;

#ifdef __cplusplus
extern "C" {
#endif

QItemModel_* newGoItemModel(QObject_*, GoAddr*);
void deleteGoItemModel(QItemModel_*);


#ifdef __cplusplus
} // extern "C"
#endif

#endif // MODELINDEX_H

// vim:ts=4:sw=4:et:ft=cpp

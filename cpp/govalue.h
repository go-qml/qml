#ifndef GOVALUE_H
#define GOVALUE_H

// Unfortunatley we need access to private bits, because the
// whole dynamic meta-object concept is sadly being hidden
// away, and without it this package wouldn't exist.
#include <private/qmetaobject_p.h>

#include "capi.h"

class GoValueMetaObject;

// TODO Painting.
#include <QQuickPaintedItem>
#include <QPainter>
#include <QtQuick/QQuickItem>
#include <QtQuick/qsgnode.h>
class GoValue : public QQuickPaintedItem
//class GoValue : public QObject
{
    Q_OBJECT

public:
    GoAddr *addr;
    GoTypeInfo *typeInfo;

    GoValue(GoAddr *addr, GoTypeInfo *typeInfo, QObject *parent);

    void activate(int propIndex);

    // TODO Painting.
    virtual void paint(QPainter *painter);
    //virtual QSGNode *updatePaintNode(QSGNode *oldNode, UpdatePaintNodeData *updatePaintNodeData);
    //virtual void itemChange(ItemChange, const ItemChangeData &){};

    static QMetaObject *metaObjectFor(GoTypeInfo *typeInfo);

    virtual ~GoValue();

private:
    GoValueMetaObject *valueMeta;

//public slots:
//    virtual void paint();
};

#endif // GOVALUE_H

// vim:ts=4:sw=4:et:ft=cpp

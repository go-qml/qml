#ifndef GOVALUE_H
#define GOVALUE_H

// Unfortunatley we need access to private bits, because the
// whole dynamic meta-object concept is sadly being hidden
// away, and without it this package wouldn't exist.
#include <private/qmetaobject_p.h>

#include <QQuickPaintedItem>
#include <QPainter>

#include "capi.h"

class GoValueMetaObject;

QMetaObject *metaObjectFor(GoTypeInfo *typeInfo);

class GoValue : public QQuickPaintedItem
{
    Q_OBJECT

public:
    GoAddr *addr;
    GoTypeInfo *typeInfo;

    GoValue(GoAddr *addr, GoTypeInfo *typeInfo, QObject *parent);

    void activate(int propIndex);

    virtual ~GoValue();

    virtual void paint(QPainter *painter);

private:
    GoValueMetaObject *valueMeta;
};

#endif // GOVALUE_H

// vim:ts=4:sw=4:et:ft=cpp


#include <private/qmetaobjectbuilder_p.h>

#include <QtQml/qqml.h>
#include <QQmlEngine>
#include <QDebug>

#include "govalue.h"
#include "capi.h"

class GoValueMetaObject : public QAbstractDynamicMetaObject
{
public:
    GoValueMetaObject(GoValue* value, GoTypeInfo *typeInfo);

protected:
    int metaCall(QMetaObject::Call c, int id, void **a);

private:
    GoValue *value;
};

GoValueMetaObject::GoValueMetaObject(GoValue *value, GoTypeInfo *typeInfo)
    : value(value)
{
    //d->parent = static_cast<QAbstractDynamicMetaObject *>(priv->metaObject);
    *static_cast<QMetaObject *>(this) = *GoValue::metaObjectFor(typeInfo);

    QObjectPrivate *objPriv = QObjectPrivate::get(value);
    objPriv->metaObject = this;
}

int GoValueMetaObject::metaCall(QMetaObject::Call c, int idx, void **a)
{
    //qWarning() << "GoValueMetaObject::metaCall" << c << idx;
    switch (c) {
    case QMetaObject::ReadProperty:
    case QMetaObject::WriteProperty:
        {
            // TODO Cache propertyOffset, methodOffset (and maybe qmlEngine)
            int propOffset = propertyOffset();
            if (idx < propOffset) {
                return value->qt_metacall(c, idx, a);
            }
            GoMemberInfo *memberInfo = value->typeInfo->fields;
            for (int i = 0; i < value->typeInfo->fieldsLen; i++) {
                if (memberInfo->metaIndex == idx) {
                    if (c == QMetaObject::ReadProperty) {
                        DataValue result;
                        hookGoValueReadField(qmlEngine(value), value->addr, memberInfo->reflectIndex, &result);
                        QVariant *out = reinterpret_cast<QVariant *>(a[0]);
                        unpackDataValue(&result, out);
                    } else {
                        DataValue assign;
                        QVariant *in = reinterpret_cast<QVariant *>(a[0]);
                        packDataValue(in, &assign);
                        hookGoValueWriteField(qmlEngine(value), value->addr, memberInfo->reflectIndex, memberInfo->reflectChangedIndex, &assign);
                        activate(value, methodOffset() + (idx - propOffset), 0);
                    }
                    return -1;
                }
                memberInfo++;
            }
            QMetaProperty prop = property(idx);
            qWarning() << "Property" << prop.name() << "not found!?";
            break;
        }
    case QMetaObject::InvokeMetaMethod:
        {
            if (idx < methodOffset()) {
                return value->qt_metacall(c, idx, a);
            }
            GoMemberInfo *memberInfo = value->typeInfo->methods;
            for (int i = 0; i < value->typeInfo->methodsLen; i++) {
                if (memberInfo->metaIndex == idx) {
                    // args[0] is the result if any.
                    DataValue args[1 + MaxParams];
                    for (int i = 1; i < memberInfo->numIn+1; i++) {
                        packDataValue(reinterpret_cast<QVariant *>(a[i]), &args[i]);
                    }
                    hookGoValueCallMethod(qmlEngine(value), value->addr, memberInfo->reflectIndex, args);
                    if (memberInfo->numOut > 0) {
                        unpackDataValue(&args[0], reinterpret_cast<QVariant *>(a[0]));
                    }
                    return -1;
                }
                memberInfo++;
            }
            QMetaMethod m = method(idx);
            qWarning() << "Method" << m.name() << "not found!?";
            break;
        }
    default:
        break; // Unhandled.
    }
    return -1;
}

GoValue::GoValue(GoAddr *addr, GoTypeInfo *typeInfo, QObject *parent)
    : addr(addr), typeInfo(typeInfo)
{
    valueMeta = new GoValueMetaObject(this, typeInfo);
    setParent(parent);

    // TODO Painting.
    //setFillColor(QColor(255,0,0));
}

GoValue::~GoValue()
{
    hookGoValueDestroyed(qmlEngine(this), addr);
}

void GoValue::activate(int propIndex)
{
    // Properties are added first, so the first fieldLen methods are in
    // fact the signals of the respective properties.
    int relativeIndex = propIndex - valueMeta->propertyOffset();
    valueMeta->activate(this, valueMeta->methodOffset() + relativeIndex, 0);
}

// TODO Painting.
//void GoValue::paint(QPainter *painter)
//{
//    qWarning() << "GoValue::paint() called";
//    painter->drawLine(10, 10, 40, 40);
//}

QMetaObject *GoValue::metaObjectFor(GoTypeInfo *typeInfo)
{
    if (typeInfo->metaObject) {
            return reinterpret_cast<QMetaObject *>(typeInfo->metaObject);
    }

    QMetaObjectBuilder mob;
    mob.setSuperClass(&QObject::staticMetaObject);
    // TODO Painting.
    //mob.setSuperClass(&QQuickPaintedItem::staticMetaObject);
    mob.setClassName(typeInfo->typeName);
    mob.setFlags(QMetaObjectBuilder::DynamicMetaObject);

    GoMemberInfo *memberInfo;
    
    memberInfo = typeInfo->fields;
    int relativePropIndex = mob.propertyCount();
    for (int i = 0; i < typeInfo->fieldsLen; i++) {
        mob.addSignal("__" + QByteArray::number(relativePropIndex) + "()");
        QMetaPropertyBuilder propb = mob.addProperty(memberInfo->memberName, "QVariant", relativePropIndex);
        propb.setWritable(true);
        memberInfo->metaIndex = relativePropIndex;
        memberInfo++;
        relativePropIndex++;
    }

    memberInfo = typeInfo->methods;
    int relativeMethodIndex = mob.methodCount();
    for (int i = 0; i < typeInfo->methodsLen; i++) {
        if (*memberInfo->resultSignature) {
            mob.addMethod(memberInfo->methodSignature, memberInfo->resultSignature);
        } else {
            mob.addMethod(memberInfo->methodSignature);
        }
        memberInfo->metaIndex = relativeMethodIndex;
        memberInfo++;
        relativeMethodIndex++;
    }

    // TODO Support default properties.
    //mob.addClassInfo("DefaultProperty", "text");

    QMetaObject *mo = mob.toMetaObject();

    // Turn the relative indexes into absolute indexes.
    memberInfo = typeInfo->fields;
    int propOffset = mo->propertyOffset();
    for (int i = 0; i < typeInfo->fieldsLen; i++) {
        memberInfo->metaIndex += propOffset;
        memberInfo++;
    }
    memberInfo = typeInfo->methods;
    int methodOffset = mo->methodOffset();
    for (int i = 0; i < typeInfo->methodsLen; i++) {
        memberInfo->metaIndex += methodOffset;
        memberInfo++;
    }

    typeInfo->metaObject = mo;
    return mo;
}

// vim:ts=4:sw=4:et:ft=cpp

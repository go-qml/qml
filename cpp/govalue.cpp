
#include <private/qmetaobjectbuilder_p.h>

#include <QQmlEngine>
#include <QtQml/qqml.h>
#include <QDebug>

#include "govalue.h"
#include "capi.h"

class GoValuePrivate;
class GoValueMetaObject : public QAbstractDynamicMetaObject
{
public:
    GoValueMetaObject(GoValue* value, GoValuePrivate *valuePriv, GoTypeInfo *typeInfo);

protected:
    int metaCall(QMetaObject::Call c, int id, void **a);

private:
    GoValue *value;
    GoValuePrivate *valuePriv;
};

class GoValuePrivate : public QObjectPrivate
{
    Q_DECLARE_PUBLIC(GoValue)
public:
    GoValueMetaObject *valueMeta;
    GoTypeInfo *typeInfo;
    GoAddr *addr;
};

GoValueMetaObject::GoValueMetaObject(GoValue *value_, GoValuePrivate *valuePriv_, GoTypeInfo *typeInfo)
    : value(value_), valuePriv(valuePriv_)
{
    //d->parent = static_cast<QAbstractDynamicMetaObject *>(priv->metaObject);
    *static_cast<QMetaObject *>(this) = *GoValue::metaObjectFor(typeInfo);

    QObjectPrivate *objPriv = QObjectPrivate::get(value);
    objPriv->metaObject = this;
}

int GoValueMetaObject::metaCall(QMetaObject::Call c, int idx, void **a)
{
    qWarning() << "metaCall:" << c << idx;
    switch (c) {
    case QMetaObject::ReadProperty:
    case QMetaObject::WriteProperty:
        {
            // TODO cache propertyOffset()
            if (idx < propertyOffset()) {
                return value->qt_metacall(c, idx, a);
            }
            GoMemberInfo *memberInfo = valuePriv->typeInfo->fields;
            for (int i = 0; i < valuePriv->typeInfo->fieldsLen; i++) {
                if (memberInfo->metaIndex == idx) {
                    // TODO Cache the qmlEngine call result?
                    if (c == QMetaObject::ReadProperty) {
                        DataValue result;
                        hookGoValueReadField(qmlEngine(value), valuePriv->addr, memberInfo->reflectIndex, &result);
                        QVariant *out = reinterpret_cast<QVariant *>(a[0]);
                        unpackDataValue(&result, out);
                    } else {
                        DataValue assign;
                        QVariant *in = reinterpret_cast<QVariant *>(a[0]);
                        packDataValue(in, &assign);
                        hookGoValueWriteField(qmlEngine(value), valuePriv->addr, memberInfo->reflectIndex, &assign);
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
            // TODO cache methodOffset()
            if (idx < methodOffset()) {
                return value->qt_metacall(c, idx, a);
            }
            GoMemberInfo *memberInfo = valuePriv->typeInfo->methods;
            for (int i = 0; i < valuePriv->typeInfo->methodsLen; i++) {
                if (memberInfo->metaIndex == idx) {
                    // TODO Cache the qmlEngine call result?
                    DataValue result;
                    hookGoValueCallMethod(qmlEngine(value), valuePriv->addr, memberInfo->reflectIndex, &result);
                    QVariant *out = reinterpret_cast<QVariant *>(a[0]);
                    unpackDataValue(&result, out);
                    return -1;
                }
                memberInfo++;
            }
            QMetaMethod m = method(idx);
            qWarning() << "Method" << m.name() << "not found!?";
            break;
        }
    }
    return -1;
}

GoValue::GoValue(GoAddr *addr, GoTypeInfo *typeInfo, QObject *parent)
        : QObject(*(new GoValuePrivate()), parent)
{
    Q_D(GoValue);
    d->addr = addr;
    d->typeInfo = typeInfo;
    d->valueMeta = new GoValueMetaObject(this, d, typeInfo);
}

GoValue::~GoValue()
{
    Q_D(GoValue);
    hookGoValueDestroyed(qmlEngine(this), d->addr);
}

GoAddr *GoValue::addr()
{
    Q_D(GoValue);
    return d->addr;
}

void GoValue::activate(int fieldReflectIndex) {
    Q_D(GoValue);

    // Go fields have an absolute index of (propertyOffset + fieldReflectIndex),
    // while Go methods have (methodOffset + fieldCount + methodReflectIndex),
    // because properties are added first, and each property gets its own signal
    // method at that time. This means the first fieldCount methods are in fact
    // the signals of the respective properties.
    d->valueMeta->activate(this, d->valueMeta->methodOffset() + fieldReflectIndex, 0);
}

QMetaObject *GoValue::metaObjectFor(GoTypeInfo *typeInfo)
{
    if (typeInfo->metaObject) {
            return reinterpret_cast<QMetaObject *>(typeInfo->metaObject);
    }

    QMetaObjectBuilder mob;
    mob.setSuperClass(&QObject::staticMetaObject);
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
        // TODO Unhardcode this.
        QMetaMethodBuilder methb = mob.addMethod("stringMethod()", "QString");
        memberInfo->metaIndex = relativeMethodIndex;
        memberInfo++;
        relativeMethodIndex++;
    }

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


#include <private/qmetaobjectbuilder_p.h>

#include <QDebug>

#include "govalue.h"


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
    Q_UNUSED(c);
    Q_UNUSED(a);
    qDebug() << "Got metaCall for" << idx << " - Reading: " << (c == QMetaObject::ReadProperty);
    if (c == QMetaObject::ReadProperty) {
        GoMemberInfo *memberInfo = valuePriv->typeInfo->members;
        for (int i = 0; i < valuePriv->typeInfo->membersLen; i++) {
            if (memberInfo->metaIndex == idx) {
                qint64 mem;
                void *result = &mem;
                hookReadField(valuePriv->addr, memberInfo->memberIndex, result);
                QVariant *out = reinterpret_cast<QVariant *>(a[0]);
                switch (memberInfo->memberType) {
                case DTString:
                    *out = *(char **)result;
                    break;
                case DTInt64:
                    *out = *(qint64 *)result;
                    break;
                case DTInt32:
                    *out = *(qint32 *)result;
                    break;
                case DTFloat64:
                    *out = *(double *)result;
                    break;
                case DTFloat32:
                    *out = *(float *)result;
                    break;
                default:
                    Q_ASSERT_X(false, "assignment", "unsupported type");
                    break;
                }
                return -1;
            }
            memberInfo++;
        }
        QMetaProperty prop = property(idx);
        qWarning() << "Property" << prop.name() << "not found!?";
    }
    return -1;
}

GoValue::GoValue(GoAddr *addr, GoTypeInfo *typeInfo)
        : QObject(*(new GoValuePrivate()), 0)
{
    Q_D(GoValue);
    d->addr = addr;
    d->typeInfo = typeInfo;
    d->valueMeta = new GoValueMetaObject(this, d, typeInfo);
}

QMetaObject *GoValue::metaObjectFor(GoTypeInfo *typeInfo)
{
    QMetaObjectBuilder mob;
    mob.setSuperClass(&QObject::staticMetaObject);
    mob.setClassName(typeInfo->typeName);
    mob.setFlags(QMetaObjectBuilder::DynamicMetaObject);

    int id = mob.propertyCount();
    GoMemberInfo *memberInfo = typeInfo->members;
    for (int i = 0; i < typeInfo->membersLen; i++) {
        mob.addSignal("__" + QByteArray::number(id) + "()");
        QMetaPropertyBuilder propb = mob.addProperty(memberInfo->memberName, "QVariant", id);
        propb.setWritable(true);
        memberInfo->metaIndex = propb.index();
        memberInfo++;
        id++;
    }

    QMetaObject *mo = mob.toMetaObject();

    int propertyOffset = mo->propertyOffset();
    memberInfo = typeInfo->members;
    for (int i = 0; i < typeInfo->membersLen; i++) {
        memberInfo->metaIndex += propertyOffset;
        memberInfo++;
    }

    // XXX Must cache mo.

    return mo;
}


// vim:ts=4:sw=4:et:ft=cpp

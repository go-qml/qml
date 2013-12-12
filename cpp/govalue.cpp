
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
                        hookGoValueReadField(qmlEngine(value), value->addr, memberInfo->reflectIndex, memberInfo->reflectChangedIndex, &result);
                        if (memberInfo->memberType == DTListProperty) {
                            if (result.dataType != DTListProperty) {
                                panicf("reading DTListProperty field returned non-DTListProperty result");
                            }
                            QQmlListProperty<QObject> *in = *reinterpret_cast<QQmlListProperty<QObject> **>(result.data);
                            QQmlListProperty<QObject> *out = reinterpret_cast<QQmlListProperty<QObject> *>(a[0]);
                            *out = *in;
                            // TODO Could provide a single variable in the stack to ReadField instead.
                            delete in;
                        } else {
                            QVariant *out = reinterpret_cast<QVariant *>(a[0]);
                            unpackDataValue(&result, out);
                        }
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

    QQuickItem::setFlag(QQuickItem::ItemHasContents, true);
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

#include <QOpenGLContext>

void GoValue::itemChange(ItemChange change, const ItemChangeData &)
{
    QQuickWindow *win = window();
    if (change != ItemSceneChange || !win) {
        return;
    }
    //QObject::connect(win, &QQuickWindow::beforeRendering, [=]() {
    //    qWarning() << "beforeRendering";
    //    glViewport(0, 0, window()->width(), window()->height());
    //    glLineWidth(2.5); 
    //    glColor3f(1.0, 0.0, 0.0);
    //    glBegin(GL_LINES);
    //    glVertex3f(0.0, 0.0, 0.0);
    //    glVertex3f(15, 0, 0);
    //    //glEnd();
    //});

    QObject::connect(win, &QQuickWindow::beforeRendering, this, &GoValue::paint, Qt::DirectConnection);
}

// TODO Painting.
//void GoValue::paint(QPainter *painter)
void GoValue::paint()
{
    qWarning() << "GoValue::paint() called";
    //painter->drawLine(10, 10, 40, 40);

    window()->setClearBeforeRendering(false);
    glViewport(0, 0, window()->width(), window()->height());
    glClear( GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT );
    glLineWidth(2.5); 
    glColor3f(1.0, 0.0, 0.0);
    glBegin(GL_LINES);
    glVertex3f(0.0, 0.0, 0.0);
    glVertex3f(1, 0, 0);
    glEnd();

    //glViewport(0, 0, 150, 150);
    //glMatrixMode(GL_PROJECTION);
    //glLoadIdentity();
    //glOrtho(0, 150, 150,0,0,10);
    //glMatrixMode(GL_MODELVIEW);
    //glLoadIdentity();

    //glClear( GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT );
    //glColor3ub(255, 0, 0);
    //glBegin(GL_QUADS);
    //    glVertex2i(0, 0);
    //    glVertex2i(100, 0);
    //    glVertex2i(100, 100);
    //    glVertex2i(0, 100);
    //glEnd();
}

QMetaObject *GoValue::metaObjectFor(GoTypeInfo *typeInfo)
{
    if (typeInfo->metaObject) {
            return reinterpret_cast<QMetaObject *>(typeInfo->metaObject);
    }

    QMetaObjectBuilder mob;
    // TODO Painting.
    mob.setSuperClass(&QQuickItem::staticMetaObject);
    //mob.setSuperClass(&QObject::staticMetaObject);
    mob.setClassName(typeInfo->typeName);
    mob.setFlags(QMetaObjectBuilder::DynamicMetaObject);

    GoMemberInfo *memberInfo;
    
    memberInfo = typeInfo->fields;
    int relativePropIndex = mob.propertyCount();
    for (int i = 0; i < typeInfo->fieldsLen; i++) {
        mob.addSignal("__" + QByteArray::number(relativePropIndex) + "()");
        const char *typeName = "QVariant";
        if (memberInfo->memberType == DTListProperty) {
            typeName = "QQmlListProperty<QObject>";
        }
        QMetaPropertyBuilder propb = mob.addProperty(memberInfo->memberName, typeName, relativePropIndex);
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
    //mob.addClassInfo("DefaultProperty", "objects");

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

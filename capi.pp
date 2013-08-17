#include <QApplication>
#include <QQuickView>
#include <QtQml>
#include <QDebug>

#include "govalue.h"
#include "govaluetype.h"
#include "capi.h"

#include <QDebug>

QApplication_ *newGuiApplication(int argc, char **argv)
{
    return new QGuiApplication(argc, argv);
}

void applicationExec(QApplication_ *app)
{
    reinterpret_cast<QCoreApplication *>(app)->exec();
}

QQmlEngine_ *newEngine(QObject_ *parent)
{
    return new QQmlEngine(reinterpret_cast<QObject *>(parent));
}

void delEngine(QQmlEngine_ *engine)
{
    QQmlEngine *qengine = reinterpret_cast<QQmlEngine *>(engine);
    delete qengine;
}

QQmlContext_ *engineRootContext(QQmlEngine_ *engine)
{
    return reinterpret_cast<QQmlEngine *>(engine)->rootContext();
}

void contextSetObject(QQmlContext_ *context, QObject_ *value)
{
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);
    QObject *qvalue = reinterpret_cast<QObject *>(value);

    qcontext->setContextObject(qvalue);
}

void contextSetPropertyObject(QQmlContext_ *context, QString_ *name, QObject_ *value)
{
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);
    const QString *qname = reinterpret_cast<QString *>(name);
    QObject *qvalue = reinterpret_cast<QObject *>(value);

    qcontext->setContextProperty(*qname, qvalue);
}

void contextSetPropertyString(QQmlContext_ *context, QString_ *name, QString_ *value)
{
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);
    const QString *qname = reinterpret_cast<QString *>(name);
    const QString *qvalue = reinterpret_cast<QString *>(value);

    qcontext->setContextProperty(*qname, *qvalue);
}

void contextSetPropertyBool(QQmlContext_ *context, QString_ *name, int32_t value)
{
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);
    const QString *qname = reinterpret_cast<QString *>(name);

    qcontext->setContextProperty(*qname, value == 0 ? false : true);
}

void contextSetPropertyInt64(QQmlContext_ *context, QString_ *name, int64_t value)
{
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);
    const QString *qname = reinterpret_cast<QString *>(name);

    qcontext->setContextProperty(*qname, qint64(value));
}

void contextSetPropertyInt32(QQmlContext_ *context, QString_ *name, int32_t value)
{
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);
    const QString *qname = reinterpret_cast<QString *>(name);

    qcontext->setContextProperty(*qname, qint32(value));
}

void contextSetPropertyFloat64(QQmlContext_ *context, QString_ *name, double value)
{
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);
    const QString *qname = reinterpret_cast<QString *>(name);

    qcontext->setContextProperty(*qname, value);
}

void contextSetPropertyFloat32(QQmlContext_ *context, QString_ *name, float value)
{
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);
    const QString *qname = reinterpret_cast<QString *>(name);

    qcontext->setContextProperty(*qname, value);
}

void contextGetProperty(QQmlContext_ *context, QString_ *name, void *result, DataType *dtype)
{
    const QString *qname = reinterpret_cast<QString *>(name);
    QVariant var = reinterpret_cast<QQmlContext *>(context)->contextProperty(*qname);

    // Some assumptions are made below regarding the size of types.
    // There's apparently no better way to handle this since that's
    // how the types with well defined sizes (qint64) are mapped to
    // meta-types (QMetaType::LongLong).
    switch ((int)var.type()) {
    case QMetaType::QString:
        {
            QByteArray ba = var.toByteArray();
            *(char**)(result) = strdup(ba.constData());
            *dtype = DTString;
            break;
        }
    case QMetaType::Bool:
        *(qint32*)(result) = var.toInt();
        *dtype = DTBool;
        break;
    case QMetaType::LongLong:
        *(qint64*)(result) = var.toLongLong();
        *dtype = DTInt64;
        break;
    case QMetaType::Int:
        *(qint32*)(result) = var.toInt();
        *dtype = DTInt32;
        break;
    case QMetaType::Double:
        *(double*)(result) = var.toDouble();
        *dtype = DTFloat64;
        break;
    case QMetaType::Float:
        *(float*)(result) = var.toFloat();
        *dtype = DTFloat32;
        break;
    case QMetaType::QObjectStar:
        {
            QObject *qobject = var.value<QObject *>();
            GoValue *value = dynamic_cast<GoValue *>(qobject);
            if (value) {
                *(void **)(result) = value->addr();
                *dtype = DTGoAddr;
                break;
            }
        }
        // fallthrough
    default:
        qWarning() << "Unsupported variant type:" << var.type();
        break;
    }
}

QString_ *newString(const char *data, int len)
{
    // This will copy data only once.
    QByteArray ba = QByteArray::fromRawData(data, len);
    return new QString(ba);
}

void delString(QString_ *s)
{
    delete reinterpret_cast<QString *>(s);
}

QObject_ *newValue(GoAddr *addr, GoTypeInfo *typeInfo)
{
    return new GoValue(addr, typeInfo);
}

int gqRunSpike(GoAddr *addr, GoTypeInfo *typeInfo)
{
    int argc = 1;
    const char *argv[] = {""};

    QGuiApplication app(argc, (char **)argv);
    QQuickView view;

    GoValue value(addr, typeInfo);
    view.rootContext()->setContextProperty("value", &value);

    GoValueType<1>::init(typeInfo);
    qmlRegisterType< GoValueType<1> >("GoTypes", 1, 0, typeInfo->typeName);

    view.setSource(QUrl::fromLocalFile("spike.qml"));
    view.show();

    return app.exec();
}

// vim:ts=4:sw=4:et:ft=cpp

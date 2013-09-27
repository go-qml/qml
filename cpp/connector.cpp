#include <QObject>

#include "connector.h"

Connector::~Connector()
{
    hookSignalDisconnect(func);
}

void Connector::invoke()
{
    qFatal("should never get called");
}

int Connector::qt_metacall(QMetaObject::Call c, int idx, void **a)
{
    if (c == QMetaObject::InvokeMetaMethod && idx == metaObject()->methodOffset()) {
        DataValue args[MaxParams];
        for (int i = 0; i < argsLen; i++) {
            QVariant var(method.parameterType(i), a[1 + i]);
            packDataValue(&var, &args[i]);
        }
        hookSignalCall(engine, func, args);
        return -1;
    }
    return standard_qt_metacall(c, idx, a);
}

// vim:ts=4:sw=4:et

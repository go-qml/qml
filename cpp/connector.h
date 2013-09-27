#ifndef CONNECTOR_H
#define CONNECTOR_H

#include <QObject>

class Connector : public QObject
{
    Q_OBJECT

    public:

    Connector(QObject *sender, QMetaMethod method, QQmlEngine *engine, void *func, int argsLen)
        : QObject(sender), engine(engine), method(method), func(func), argsLen(argsLen) {};

    virtual ~Connector();

    // MOC HACK: s/Connector::qt_metacall/Connector::standard_qt_metacall/
    int standard_qt_metacall(QMetaObject::Call c, int idx, void **a);

    public slots:

    void invoke();

    private:

    QQmlEngine *engine;
    QMetaMethod method;
    void *func;
    int argsLen;
};

#endif // CONNECTOR_H

// vim:ts=4:sw=4:et

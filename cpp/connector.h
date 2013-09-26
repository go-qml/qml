#ifndef CONNECTOR_H
#define CONNECTOR_H

#include <QObject>

class Connector : public QObject
{
    Q_OBJECT

    public:

    Connector(QQmlEngine *engine, QObject *sender, QMetaMethod method, void *data, int argsLen)
        : engine(engine), sender(sender), method(method), data(data), argsLen(argsLen) {};

    // MOC HACK: s/Connector::qt_metacall/Connector::standard_qt_metacall/
    int standard_qt_metacall(QMetaObject::Call c, int idx, void **a);

    public slots:

    void invoke();

    private:

    QQmlEngine *engine;
    QObject *sender;
    QMetaMethod method;
    void *data;
    int argsLen;
};

#endif // CONNECTOR_H

// vim:ts=4:sw=4:et

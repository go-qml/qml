#ifndef CONNECTOR_H
#define CONNECTOR_H

#include <QObject>

class Connector : public QObject
{
    Q_OBJECT

    public:

    Connector(QObject *sender, void *data) : sender(sender), data(data) {};

    public slots:

    void invoke0();

    private:

    QObject *sender;
    void *data;
};

#endif // CONNECTOR_H

// vim:ts=4:sw=4:et

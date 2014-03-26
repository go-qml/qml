#ifndef TESTTYPE_H
#define TESTTYPE_H

#include <QObject>

class TestType : public QObject
{
    Q_OBJECT

    Q_PROPERTY(void *voidAddr READ getVoidAddr)

    void *voidAddr;

    public:

    TestType(QObject *parent = 0) : QObject(parent), voidAddr((void*)42) {};

    void *getVoidAddr() { return voidAddr; };
};

#endif // TESTTYPE_H

// vim:ts=4:sw=4:et

#ifndef IDLE_TIMER_H
#define IDLE_TIMER_H

#include <QBasicTimer>
#include <QThread>
#include <QDebug>
#include <mutex>

#include "goqml_private.h"

class IdleTimer : public QObject
{
    Q_OBJECT

    public:

    static IdleTimer *singleton() {
        static IdleTimer singleton;
        return &singleton;
    }

    void init(int32_t *guiIdleRun)
    {
        this->guiIdleRun = guiIdleRun;
    }

    Q_INVOKABLE void start()
    {
        timer.start(0, this);
    }

    protected:

    void timerEvent(QTimerEvent *event);

    private:

    int32_t *guiIdleRun;

    QBasicTimer timer;    
};

#endif

// vim:ts=4:sw=4:et:ft=cpp

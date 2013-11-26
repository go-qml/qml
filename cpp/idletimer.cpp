#include <QBasicTimer>
#include <QThread>
#include <QDebug>
#include <mutex>

#include "capi.h"

class IdleTimer : public QObject
{
    Q_OBJECT

    public:

    static IdleTimer *singleton() {
        static IdleTimer singleton;
        return &singleton;
    }

    void init(int *hookWaiting)
    {
        this->hookWaiting = hookWaiting;
    }

    Q_INVOKABLE void start()
    {
        timer.start(0, this);
    }

    protected:

    void timerEvent(QTimerEvent *event)
    {
        __sync_synchronize();
        if (*hookWaiting > 0) {
            hookIdleTimer();
        } else {
            timer.stop();
        }
    }

    private:

    int *hookWaiting;

    QBasicTimer timer;    
};

void idleTimerInit(int *hookWaiting)
{
    IdleTimer::singleton()->init(hookWaiting);
}

void idleTimerStart()
{
    QMetaObject::invokeMethod(IdleTimer::singleton(), "start", Qt::QueuedConnection);
}

// vim:ts=4:sw=4:et:ft=cpp

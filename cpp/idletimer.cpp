#include <QBasicTimer>
#include <QThread>
#include <QDebug>

#include "capi.h"

extern "C" {

int g_atomic_int_get(const volatile int *value);

}

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
        if (g_atomic_int_get(hookWaiting) > 0) {
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

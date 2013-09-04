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

    static void start(int *hookWaiting)
    {
        static IdleTimer singleton;
        singleton.hookWaiting = hookWaiting;
        singleton.timer.start(0, &singleton);
    }

    protected:

    void timerEvent(QTimerEvent *event)
    {
        if (g_atomic_int_get(hookWaiting) > 0) {
            hookIdleTimer();
        }
    }

    private:

    int *hookWaiting;

    QBasicTimer timer;    
};

void startIdleTimer(int *hookWaiting)
{
    IdleTimer::start(hookWaiting);
}

// vim:ts=4:sw=4:et:ft=cpp

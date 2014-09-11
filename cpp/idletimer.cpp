#include "idletimer.h"
#include "goqml_private.h"


#if defined(_MSC_VER)
#   include <Windows.h>
#endif

void idleTimerInit(int32_t *guiIdleRun)
{
    IdleTimer::singleton()->init(guiIdleRun);
}

void idleTimerStart()
{
    QMetaObject::invokeMethod(IdleTimer::singleton(), "start", Qt::QueuedConnection);
}

    void IdleTimer::timerEvent(QTimerEvent *event)
    {
#if defined(_MSC_VER)
        MemoryBarrier();
#else
        __sync_synchronize();
#endif
        if (*guiIdleRun > 0) {
            hookIdleTimer();
        } else {
            timer.stop();
        }
    }

// vim:ts=4:sw=4:et:ft=cpp

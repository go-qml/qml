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
        if (_singleton == nullptr) 
            _singleton = new IdleTimer;
        return _singleton;
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

    void timerEvent(QTimerEvent *event)
    {
        __sync_synchronize();
        if (*guiIdleRun > 0) {
            hookIdleTimer();
        } else {
            timer.stop();
        }
    }

    static IdleTimer* _singleton;

    private:

    int32_t *guiIdleRun;

    QBasicTimer timer;    
};

IdleTimer* IdleTimer::_singleton;

void idleTimerInit(int32_t *guiIdleRun)
{
    IdleTimer::singleton()->init(guiIdleRun);
}

void idleTimerStart()
{
    QMetaObject::invokeMethod(IdleTimer::singleton(), "start", Qt::QueuedConnection);
}

// vim:ts=4:sw=4:et:ft=cpp

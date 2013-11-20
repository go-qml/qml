#include <QBasicTimer>
#include <QThread>
#include <QDebug>

#include "capi.h"
#include "mutex"

class IdleTimer : public QObject
{
    Q_OBJECT

    public:

    static IdleTimer *singleton() {
        if (!instance) {
            
            std::lock_guard<std::mutex> lock(mx);

            if (!instance) {
                instance = new IdleTimer();
            }
        }
        return instance;
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
        // this is a gcc intrinsic, non-gcc compilers will
        // need some other memory barrier inducing call
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

    static std::mutex mx;
    static IdleTimer *instance;
};

IdleTimer* IdleTimer::instance = 0;
std::mutex IdleTimer::mx;

void idleTimerInit(int *hookWaiting)
{
    IdleTimer::singleton()->init(hookWaiting);
}

void idleTimerStart()
{
    QMetaObject::invokeMethod(IdleTimer::singleton(), "start", Qt::QueuedConnection);
}

// vim:ts=4:sw=4:et:ft=cpp

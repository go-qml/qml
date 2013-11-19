#include <QBasicTimer>
#include <QThread>
#include <QDebug>

#include "capi.h"

#ifdef __MINGW64__
#include "mutex"
#else
extern "C" {

int g_atomic_int_get(const volatile int *value);

}
#endif 

class IdleTimer : public QObject
{
    Q_OBJECT

    public:

    static IdleTimer *singleton() {
#ifdef __MINGW64__
        if (!instance) {
            
            std::lock_guard<std::mutex> lock(mx);

            if (!instance) {
                instance = new IdleTimer();
            }
        }
        return instance;
#else
        static IdleTimer singleton;
        return &singleton;
#endif
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
#ifdef __MINGW64__  
        MemoryBarrier();

        if (*hookWaiting > 0) {  
#else    
        if (g_atomic_int_get(hookWaiting) > 0) {
#endif
            hookIdleTimer();
        } else {
            timer.stop();
        }
    }

    private:

    int *hookWaiting;
    QBasicTimer timer;

#ifdef __MINGW64__
    static std::mutex mx;
    static IdleTimer *instance;
#endif
};

#ifdef __MINGW64__
IdleTimer* IdleTimer::instance = 0;
std::mutex IdleTimer::mx;
#endif 

void idleTimerInit(int *hookWaiting)
{
    IdleTimer::singleton()->init(hookWaiting);
}

void idleTimerStart()
{
    QMetaObject::invokeMethod(IdleTimer::singleton(), "start", Qt::QueuedConnection);
}

// vim:ts=4:sw=4:et:ft=cpp

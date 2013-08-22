#include <QBasicTimer>
#include <QThread>
#include <QDebug>

#include "capi.h"

class IdleTimer : public QObject
{
    Q_OBJECT

    public:

    static void start()
    {
        static IdleTimer singleton;
        singleton.timer.start(0, &singleton);
    }

    protected:

    int c;

    void timerEvent(QTimerEvent *event)
    {
        if (c == 100) {
            qDebug() << "IdleTimer::timerEvent is in" << QThread::currentThread();
        }
        c++;
        // Might be worth sharing some synchronized flag to tell
        // whether there's work to do or not.
        hookIdleTimer();
    }

    private:

    QBasicTimer timer;    
};

void startIdleTimer()
{
    IdleTimer::start();
}

// vim:ts=4:sw=4:et:ft=cpp

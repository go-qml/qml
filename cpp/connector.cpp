#include <QObject>

#include "connector.h"


void Connector::invoke0()
{
    hookSignal(sender, data, 0, 0);
}

// vim:ts=4:sw=4:et

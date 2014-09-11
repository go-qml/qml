#ifndef GOQML_PRIVATE_H_
#define GOQML_PRIVATE_H_

#include "goqml.h"

#ifdef __cplusplus
extern "C" {
#endif

void hookIdleTimer();
void hookLogHandler(LogMessage *message);
void hookGoValueReadField(QQmlEngine_ *engine, GoAddr *addr, int memberIndex, int getIndex, int setIndex, DataValue *result);
void hookGoValueWriteField(QQmlEngine_ *engine, GoAddr *addr, int memberIndex, int setIndex, DataValue *assign);
void hookGoValueCallMethod(QQmlEngine_ *engine, GoAddr *addr, int memberIndex, DataValue *result);
void hookGoValueDestroyed(QQmlEngine_ *engine, GoAddr *addr);
void hookGoValuePaint(QQmlEngine_ *engine, GoAddr *addr, intptr_t reflextIndex);
QImage_ *hookRequestImage(void *imageFunc, char *id, int idLen, int width, int height);
GoAddr *hookGoValueTypeNew(GoValue_ *value, GoTypeSpec_ *spec);
void hookWindowHidden(QObject_ *addr);
void hookSignalCall(QQmlEngine_ *engine, void *func, DataValue *params);
void hookSignalDisconnect(void *func);
void hookPanic(char *message);
int hookListPropertyCount(GoAddr *addr, intptr_t reflectIndex, intptr_t setIndex);
QObject_ *hookListPropertyAt(GoAddr *addr, intptr_t reflectIndex, intptr_t setIndex, int i);
void hookListPropertyAppend(GoAddr *addr, intptr_t reflectIndex, intptr_t setIndex, QObject_ *obj);
void hookListPropertyClear(GoAddr *addr, intptr_t reflectIndex, intptr_t setIndex);

#ifdef __cplusplus
}
#endif
#endif	// GOQML_PRIVATE_H_


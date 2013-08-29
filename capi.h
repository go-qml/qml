#ifndef CAPI_H
#define CAPI_H

#include <stdint.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef void GoAddr;

typedef void QApplication_;
typedef void QObject_;
typedef void QVariant_;
typedef void QString_;
typedef void QQmlEngine_;
typedef void QQmlContext_;
typedef void QQmlComponent_;
typedef void QQuickView_;
typedef void QMessageLogContext_;
typedef void GoValue_;

typedef enum {
    DTUnknown = 0, // Has an unsupported type.
    DTInvalid = 1, // Does not exist or similar.
    DTAny     = 2, // Any of the following types. Used in type information, not in an actual DataValue.

    DTString  = 10,
    DTBool    = 11,
    DTInt64   = 12,
    DTInt32   = 13,
    DTFloat64 = 14,
    DTFloat32 = 15,

    DTGoAddr  = 100,
    DTObject  = 101
} DataType;

typedef struct {
    DataType dataType;
    char data[8];
    int len;
} DataValue;

typedef struct {
    char *memberName; // points to memberNames
    DataType memberType;
    int memberIndex;
    int metaIndex;
} GoMemberInfo;

typedef struct {
    char *typeName;
    GoMemberInfo *members;
    char *memberNames;
    int membersLen;
} GoTypeInfo;

typedef struct {
    int severity;
    const char *text;
    int textLen;
    const char *file;
    int fileLen;
    int line;
} LogMessage;

void newGuiApplication();
void applicationExec();
void applicationFlushAll();
void startIdleTimer();

void *currentThread();
void *appThread();

QQmlEngine_ *newEngine(QObject_ *parent);
QQmlContext_ *engineRootContext(QQmlEngine_ *engine);
void engineSetOwnershipCPP(QQmlEngine_ *engine, QObject_ *object);
void engineSetOwnershipJS(QQmlEngine_ *engine, QObject_ *object);
void engineSetContextForObject(QQmlEngine_ *engine, QObject_ *object);

void contextGetProperty(QQmlContext_ *context, QString_ *name, DataValue *value);
void contextSetProperty(QQmlContext_ *context, QString_ *name, DataValue *value);
void contextSetObject(QQmlContext_ *context, QObject_ *value);

void delObject(QObject_ *object);
void objectGetProperty(QObject_ *object, const char *name, DataValue *value);
void objectSetParent(QObject_ *object, QObject_ *parent);

QQmlComponent_ *newComponent(QQmlEngine_ *engine, QObject_ *parent);
void componentSetData(QQmlComponent_ *component, const char *data, int dataLen, const char *url, int urlLen);
char *componentErrorString(QQmlComponent_ *component);
QObject_ *componentCreate(QQmlComponent_ *component, QQmlContext_ *context);
QQuickView_ *componentCreateView(QQmlComponent_ *component, QQmlContext_ *context);

void viewShow(QQuickView_ *view);
void viewHide(QQuickView_ *view);
void viewReportHidden(QQuickView_ *view);
QObject_ *viewRootObject(QQuickView_ *view);

QString_ *newString(const char *data, int len);
void delString(QString_ *s);

GoValue_ *newValue(GoAddr *addr, GoTypeInfo *typeInfo, QObject_ *parent);

void packDataValue(QVariant_ *var, DataValue *result);
void unpackDataValue(DataValue *value, QVariant_ *result);

void installLogHandler();

void hookIdleTimer();
void hookLogHandler(LogMessage *message);
void hookGoValueReadField(QQmlEngine_ *engine, GoAddr *addr, int memberIndex, DataValue *result);
void hookGoValueDestroyed(QQmlEngine_ *engine, GoAddr *addr);
void hookWindowHidden(QObject_ *addr);

#ifdef __cplusplus
} // extern "C"
#endif

#endif // CAPI_H

// vim:ts=4:et

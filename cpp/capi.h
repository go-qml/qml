#ifndef CAPI_H
#define CAPI_H

#include <stdint.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

// It's surprising that MaximumParamCount is privately defined within qmetaobject.cpp.
// Must fix the objectInvoke function if this is changed.
// This is Qt's MaximuParamCount - 1, as it does not take the result value in account.
enum { MaxParams = 10 };

typedef void QApplication_;
typedef void QMetaObject_;
typedef void QObject_;
typedef void QVariant_;
typedef void QVariantList_;
typedef void QString_;
typedef void QQmlEngine_;
typedef void QQmlContext_;
typedef void QQmlComponent_;
typedef void QQuickWindow_;
typedef void QQuickView_;
typedef void QMessageLogContext_;
typedef void QImage_;
typedef void GoValue_;
typedef void GoAddr;
typedef void GoTypeSpec_;

typedef char error;
error *errorf(const char *format, ...);
void panicf(const char *format, ...);

typedef enum {
    DTUnknown = 0, // Has an unsupported type.
    DTInvalid = 1, // Does not exist or similar.

    DTString  = 10,
    DTBool    = 11,
    DTInt64   = 12,
    DTInt32   = 13,
    DTFloat64 = 14,
    DTFloat32 = 15,
    DTColor   = 16,

    DTGoAddr  = 100,
    DTObject  = 101,
    DTList    = 102,

    // Used in type information, not in an actual data value.
    DTAny     = 201, // Can hold any of the above types.
    DTMethod  = 202
} DataType;

typedef struct {
    DataType dataType;
    char data[8];
    int len;
} DataValue;

typedef struct {
    char *memberName; // points to memberNames
    DataType memberType;
    int reflectIndex;
    int metaIndex;
    int addrOffset;
    char *methodSignature;
    char *resultSignature;
    int numIn;
    int numOut;
} GoMemberInfo;

typedef struct {
    char *typeName;
    GoMemberInfo *fields;
    GoMemberInfo *methods;
    GoMemberInfo *members; // fields + methods
    int fieldsLen;
    int methodsLen;
    int membersLen;
    char *memberNames;

    QMetaObject_ *metaObject;
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

void idleTimerInit(int *hookWaiting);
void idleTimerStart();

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
void delObjectLater(QObject_ *object);
int objectGetProperty(QObject_ *object, const char *name, DataValue *result);
void objectSetProperty(QObject_ *object, const char *name, DataValue *value);
void objectSetParent(QObject_ *object, QObject_ *parent);
error *objectInvoke(QObject_ *object, const char *method, DataValue *result, DataValue *params, int paramsLen);
void objectFindChild(QObject_ *object, QString_ *name, DataValue *result);
QQmlContext_ *objectContext(QObject_ *object);
int objectIsComponent(QObject_ *object);
int objectIsWindow(QObject_ *object);
int objectIsView(QObject_ *object);
error *objectConnect(QObject_ *object, const char *signal, int signalLen, QQmlEngine_ *engine, void *func, int argsLen);

QQmlComponent_ *newComponent(QQmlEngine_ *engine, QObject_ *parent);
void componentSetData(QQmlComponent_ *component, const char *data, int dataLen, const char *url, int urlLen);
char *componentErrorString(QQmlComponent_ *component);
QObject_ *componentCreate(QQmlComponent_ *component, QQmlContext_ *context);
QQuickWindow_ *componentCreateWindow(QQmlComponent_ *component, QQmlContext_ *context);

void windowShow(QQuickWindow_ *win);
void windowHide(QQuickWindow_ *win);
void windowConnectHidden(QQuickWindow_ *win);
QObject_ *windowRootObject(QQuickWindow_ *win);
QImage_ *windowGrabWindow(QQuickWindow_ *win);

void delImage(QImage_ *image);
void imageSize(QImage_ *image, int *width, int *height);
const unsigned char *imageBits(QImage_ *image);

QString_ *newString(const char *data, int len);
void delString(QString_ *s);

GoValue_ *newGoValue(GoAddr *addr, GoTypeInfo *typeInfo, QObject_ *parent);
void goValueActivate(GoValue_ *value, GoTypeInfo *typeInfo, int addrOffset);

void packDataValue(QVariant_ *var, DataValue *result);
void unpackDataValue(DataValue *value, QVariant_ *result);

QVariantList_ *newVariantList(DataValue *list, int len);

int registerType(char *location, int major, int minor, char *name, GoTypeInfo *typeInfo, GoTypeSpec_ *spec);
int registerSingleton(char *location, int major, int minor, char *name, GoTypeInfo *typeInfo, GoTypeSpec_ *spec);

void installLogHandler();

void hookIdleTimer();
void hookLogHandler(LogMessage *message);
void hookGoValueReadField(QQmlEngine_ *engine, GoAddr *addr, int memberIndex, DataValue *result);
void hookGoValueWriteField(QQmlEngine_ *engine, GoAddr *addr, int memberIndex, DataValue *assign);
void hookGoValueCallMethod(QQmlEngine_ *engine, GoAddr *addr, int memberIndex, DataValue *result);
void hookGoValueDestroyed(QQmlEngine_ *engine, GoAddr *addr);
GoAddr *hookGoValueTypeNew(GoValue_ *value, GoTypeSpec_ *spec);
void hookWindowHidden(QObject_ *addr);
void hookSignalCall(QQmlEngine_ *engine, void *func, DataValue *params);
void hookSignalDisconnect(void *func);
void hookPanic(char *message);

#ifdef __cplusplus
} // extern "C"
#endif

#endif // CAPI_H

// vim:ts=4:et

#ifndef CAPI_H
#define CAPI_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef void GoAddr;

typedef void QApplication_;
typedef void QObject_;
typedef void QString_;
typedef void QQmlEngine_;
typedef void QQmlContext_;
typedef void QQmlComponent_;
typedef void QMessageLogContext_;
typedef void GoValue_;

typedef enum {
    DTString  = 1,
    DTBool    = 2,
    DTInt64   = 3,
    DTInt32   = 4,
    DTFloat64 = 5,
    DTFloat32 = 6,
    DTGoAddr  = 100,
} DataType;

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

QApplication_ *newGuiApplication(int argc, char **argv);
void applicationExec(QApplication_ *app);

QQmlEngine_ *newEngine(QObject_ *parent);
void delEngine(QQmlEngine_ *engine);
QQmlContext_ *engineRootContext(QQmlEngine_ *engine);

void contextGetProperty(QQmlContext_ *context, QString_ *name, void *result, DataType *type);
void contextSetPropertyObject(QQmlContext_ *context, QString_ *name, QObject_ *value);
void contextSetPropertyString(QQmlContext_ *context, QString_ *name, QString_ *value);
void contextSetPropertyBool(QQmlContext_ *context, QString_ *name, int32_t value);
void contextSetPropertyInt64(QQmlContext_ *context, QString_ *name, int64_t value);
void contextSetPropertyInt32(QQmlContext_ *context, QString_ *name, int32_t value);
void contextSetPropertyFloat64(QQmlContext_ *context, QString_ *name, double value);
void contextSetPropertyFloat32(QQmlContext_ *context, QString_ *name, float value);
void contextSetObject(QQmlContext_ *context, QObject_ *value);

QQmlComponent_ *newComponent(QQmlEngine_ *engine, QObject_ *parent);
QObject_ *componentCreate(QQmlComponent_ *component, QQmlContext_ *context);
void componentSetData(QQmlComponent_ *component, const char *data, int dataLen, const char *url, int urlLen);
char *componentErrorString(QQmlComponent_ *component);

QString_ *newString(const char *data, int len);
void delString(QString_ *s);

GoValue_ *newValue(GoAddr *addr, GoTypeInfo *typeInfo);

void installLogHandler();

void hookReadField(GoAddr *addr, int memberIndex, void *result);
void hookLogHandler(LogMessage *message);

#ifdef __cplusplus
} // extern "C"
#endif

#endif // CAPI_H

// vim:ts=4:et

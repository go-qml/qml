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

typedef enum {
    DTString = 1,
    DTInt64 = 2,
    DTInt32 = 3,
    DTFloat64 = 4,
    DTFloat32 = 5,
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

QApplication_ *newGuiApplication(int argc, char **argv);
void applicationExec(QApplication_ *app);

QQmlEngine_ *newEngine(QObject_ *parent);
QQmlContext_ *engineRootContext(QQmlEngine_ *engine);

void contextGet(QQmlContext_ *context, QString_ *name, void *result, DataType *type);
void contextSetObject(QQmlContext_ *context, QString_ *name, QObject_ *value);
void contextSetString(QQmlContext_ *context, QString_ *name, QString_ *value);
void contextSetInt64(QQmlContext_ *context, QString_ *name, int64_t value);
void contextSetInt32(QQmlContext_ *context, QString_ *name, int32_t value);
void contextSetFloat64(QQmlContext_ *context, QString_ *name, double value);
void contextSetFloat32(QQmlContext_ *context, QString_ *name, float value);

QString_ *newString(const char *data, int len);
void delString(QString_ *s);

QObject_ *newValue(GoAddr *addr, GoTypeInfo *typeInfo);

void hookReadField(GoAddr *addr, int memberIndex, void *result);

#ifdef __cplusplus
} // extern "C"
#endif

#endif // CAPI_H

// vim:ts=4:et

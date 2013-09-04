#ifndef GOVALUETYPE_H
#define GOVALUETYPE_H

#include "govalue.h"

template <int N>
class GoValueType : public GoValue
{
public:

    GoValueType()
        : GoValue(hookGoValueTypeNew(this, typeData), typeInfo, 0) {};

    static void init(GoTypeInfo *info, void *data)
    {
        typeInfo = info;
        typeData = data;
        static_cast<QMetaObject &>(staticMetaObject) = *GoValue::metaObjectFor(typeInfo);
    };

    static void *typeData;
    static GoTypeInfo *typeInfo;
    static QMetaObject staticMetaObject;
};

#endif // GOVALUETYPE_H

// vim:ts=4:sw=4:et

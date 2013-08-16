#ifndef GOVALUETYPE_H
#define GOVALUETYPE_H

#include "govalue.h"

template <int N>
class GoValueType : public GoValue
{
public:
    GoValueType() : GoValue(0, typeInfo) {};

    static void init(GoTypeInfo *info)
    {
        typeInfo = info;
        static_cast<QMetaObject &>(staticMetaObject) = *GoValueType::metaObjectFor(typeInfo);
    };

    static GoTypeInfo *typeInfo;
    static QMetaObject staticMetaObject;
};

#endif // GOVALUETYPE_H

// vim:ts=4:sw=4:et

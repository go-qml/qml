#include "govaluetype.h"

#define DEFINE_GOVALUETYPE(N) \
    template<> QMetaObject GoValueType<N>::staticMetaObject = QMetaObject(); \
    template<> GoTypeInfo *GoValueType<N>::typeInfo = 0; \
    template<> GoTypeSpec_ *GoValueType<N>::typeSpec = 0;

DEFINE_GOVALUETYPE(1)

// TODO Define N of these.

// vim:sw=4:st=4:et:ft=cpp

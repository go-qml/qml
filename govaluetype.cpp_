#include "govaluetype.h"

#define DEFINE_GOVALUETYPE(N) \
    template<> QMetaObject GoValueType<N>::staticMetaObject = QMetaObject(); \
    template<> GoTypeInfo *GoValueType<N>::typeInfo = 0;

DEFINE_GOVALUETYPE(1);

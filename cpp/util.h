#ifndef UTIL_H
#define UTIL_H

#include <QString>


#ifdef __cplusplus
extern "C" {
#endif



QString QStringFromGoString(const char *str, unsigned int length);



#ifdef __cplusplus
} // extern "C"
#endif

#endif // UTIL_H

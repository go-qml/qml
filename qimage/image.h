#ifndef IMAGE_H
#define IMAGE_H

#include "capi.h"

#ifdef __cplusplus
extern "C" {
#endif

QImage_* newQImage(int width, int height, unsigned int format);
QImage_* loadQImage(const char *filename, int filename_length, const char *format);
void deleteQImage(QImage_*);


#ifdef __cplusplus
} // extern "C"
#endif

#endif // IMAGE_H

// vim:ts=4:sw=4:et:ft=cpp

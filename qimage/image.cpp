
#include <QImage>
#include "image.h"
#include "util.cpp"

QImage_* newQImage(int width, int height, unsigned int format) {
  return new QImage(width, height, (QImage::Format)format);
}

QImage_* loadQImage(const char *filename, int filename_length, const char *format) {
  QString qfilename = QStringFromGoString(filename, filename_length);
  return new QImage(qfilename, format);
}

void deleteQImage(QImage_* img) {
  delete (QImage*)img;
}

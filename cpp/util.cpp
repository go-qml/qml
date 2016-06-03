
#include "util.h"

QString QStringFromGoString(const char *str, unsigned int length) {
  QByteArray qstrarray(str, length);
  return QString::fromUtf8(qstrarray);
}

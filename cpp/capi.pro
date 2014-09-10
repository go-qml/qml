TEMPLATE = lib
CONFIG  += dll release
CONFIG  -= embed_manifest_exe embed_manifest_dll
QT      += gui widgets qml quick
TARGET   = go-qml

DESTDIR = $${PWD}
INCLUDEPATH += .

HEADERS += ./govalue.h
SOURCES += ./govalue.cpp
HEADERS += ./govaluetype.h
SOURCES += ./govaluetype.cpp

# MOC HACK
# HEADERS += ./connector.h
SOURCES += ./connector.cpp ./moc_connector.cpp

HEADERS += ./idletimer.h
SOURCES += ./idletimer.cpp

HEADERS += ./capi.h
SOURCES += ./capi.cpp
DEF_FILE+= ./capi.def

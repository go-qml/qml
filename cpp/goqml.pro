TEMPLATE = lib
CONFIG  += dll release
CONFIG  -= embed_manifest_exe embed_manifest_dll
QT      += gui widgets qml quick
TARGET   = goqml

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

HEADERS += ./goqml.h
HEADERS += ./goqml_private.h
SOURCES += ./goqml.cpp
DEF_FILE+= ./goqml.def


TEMPLATE = lib
CONFIG  += dll release
CONFIG  -= embed_manifest_exe embed_manifest_dll
QT      += opengl gui
TARGET   = goqgl_3_2core

DESTDIR = $${PWD}
INCLUDEPATH += ..

HEADERS += ../funcs.h
SOURCES += ../funcs.cpp

DEF_FILE+= ./goqgl.def

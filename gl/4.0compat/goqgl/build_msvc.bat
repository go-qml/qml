
@echo off

cd %~dp0
setlocal

:: NMake: goqgl_4_0compat.dll
qmake
nmake clean
nmake release

:: MinGW: generate libgoqgl_4_0compat.a
dlltool -dllname goqgl_4_0compat.dll --def goqgl.def --output-lib libgoqgl_4_0compat.a

:: install
copy goqgl_4_0compat.dll %QTDIR%\bin


@echo off

cd %~dp0
setlocal

:: NMake: goqgl_3_3compat.dll
qmake
nmake clean
nmake release

:: MinGW: generate libgoqgl_3_3compat.a
dlltool -dllname goqgl_3_3compat.dll --def goqgl.def --output-lib libgoqgl_3_3compat.a

:: install
copy goqgl_3_3compat.dll %QTDIR%\bin

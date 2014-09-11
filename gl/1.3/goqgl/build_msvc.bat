
@echo off

cd %~dp0
setlocal

:: NMake: goqgl_1_3.dll
qmake
nmake clean
nmake release

:: MinGW: generate libgoqgl_1_3.a
dlltool -dllname goqgl_1_3.dll --def goqgl.def --output-lib libgoqgl_1_3.a

:: install
copy goqgl_1_3.dll %QTDIR%\bin

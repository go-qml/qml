
@echo off

cd %~dp0
setlocal

:: NMake: goqgl_1_4.dll
qmake
nmake clean
nmake release

:: MinGW: generate libgoqgl_1_4.a
dlltool -dllname goqgl_1_4.dll --def goqgl.def --output-lib libgoqgl_1_4.a

:: install
copy goqgl_1_4.dll %QTDIR%\bin

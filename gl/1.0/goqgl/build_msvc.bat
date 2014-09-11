
@echo off

cd %~dp0
setlocal

:: NMake: goqgl_1_0.dll
qmake
nmake clean
nmake release

:: MinGW: generate libgoqgl_1_0.a
dlltool -dllname goqgl_1_0.dll --def goqgl.def --output-lib libgoqgl_1_0.a

:: install
copy goqgl_1_0.dll %QTDIR%\bin


@echo off

cd %~dp0
setlocal

:: NMake: goqgl_3_0.dll
qmake
nmake clean
nmake release

:: MinGW: generate libgoqgl_3_0.a
dlltool -dllname goqgl_3_0.dll --def goqgl.def --output-lib libgoqgl_3_0.a

:: install
copy goqgl_3_0.dll %QTDIR%\bin

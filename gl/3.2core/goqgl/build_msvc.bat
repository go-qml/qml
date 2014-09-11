
@echo off

cd %~dp0
setlocal

:: NMake: goqgl_3_2core.dll
qmake
nmake clean
nmake release

:: MinGW: generate libgoqgl_3_2core.a
dlltool -dllname goqgl_3_2core.dll --def goqgl.def --output-lib libgoqgl_3_2core.a

:: install
copy goqgl_3_2core.dll %QTDIR%\bin

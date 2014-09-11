
@echo off

cd %~dp0
setlocal

:: NMake: goqgl_4_2core.dll
qmake
nmake clean
nmake release

:: MinGW: generate libgoqgl_4_2core.a
dlltool -dllname goqgl_4_2core.dll --def goqgl.def --output-lib libgoqgl_4_2core.a

:: install
copy goqgl_4_2core.dll %QTDIR%\bin

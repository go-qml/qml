
@echo off

cd %~dp0
setlocal

:: NMake: goqgl_4_0core.dll
qmake
nmake clean
nmake release

:: MinGW: generate libgoqgl_4_0core.a
dlltool -dllname goqgl_4_0core.dll --def goqgl.def --output-lib libgoqgl_4_0core.a

:: install
copy goqgl_4_0core.dll %QTDIR%\bin

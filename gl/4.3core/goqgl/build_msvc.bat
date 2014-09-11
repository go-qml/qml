
@echo off

cd %~dp0
setlocal

:: NMake: goqgl_4_3core.dll
qmake
nmake clean
nmake release

:: MinGW: generate libgoqgl_4_3core.a
dlltool -dllname goqgl_4_3core.dll --def goqgl.def --output-lib libgoqgl_4_3core.a

:: install
copy goqgl_4_3core.dll %QTDIR%\bin

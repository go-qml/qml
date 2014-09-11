
@echo off

cd %~dp0
setlocal

:: NMake: goqgl_es2.dll
qmake
nmake clean
nmake release

:: MinGW: generate libgoqgl_es2.a
dlltool -dllname goqgl_es2.dll --def goqgl.def --output-lib libgoqgl_es2.a

:: install
copy goqgl_es2.dll %QTDIR%\bin

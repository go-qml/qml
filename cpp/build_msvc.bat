:: Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
:: Use of this source code is governed by a BSD-style
:: license that can be found in the LICENSE file.

@echo off

cd %~dp0
setlocal

:: ----------------------------------------------------------------------------
:: NMake: goqml.dll

qmake
nmake clean
nmake release

:: ----------------------------------------------------------------------------
:: MinGW: generate libgoqml.a

dlltool -dllname goqml.dll --def goqml.def --output-lib libgoqml.a
copy goqml.dll %QTDIR%\bin

:: ----------------------------------------------------------------------------
:: PAUSE


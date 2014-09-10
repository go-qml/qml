:: Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
:: Use of this source code is governed by a BSD-style
:: license that can be found in the LICENSE file.

:: ----------------------------------------------------------------------------
:: NMake: goqgl.dll

qmake
nmake clean
nmake release

:: ----------------------------------------------------------------------------
:: MinGW: generate libgoqgl.a

dlltool -dllname goqgl.dll --def goqgl.def --output-lib libgoqgl.a
copy goqgl.dll %QTDIR%\bin

:: ----------------------------------------------------------------------------
:: PAUSE


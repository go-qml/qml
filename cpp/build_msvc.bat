:: Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
:: Use of this source code is governed by a BSD-style
:: license that can be found in the LICENSE file.

:: ----------------------------------------------------------------------------
:: NMake: go-qml.dll

qmake
nmake clean
nmake release

:: ----------------------------------------------------------------------------
:: MinGW: generate libgoqml.a

dlltool -dllname goqml.dll --def goqml.def --output-lib libgoqml.a
copy goqml.dll ..

:: ----------------------------------------------------------------------------
:: PAUSE


:: Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
:: Use of this source code is governed by a BSD-style
:: license that can be found in the LICENSE file.

:: ----------------------------------------------------------------------------
:: NMake: go-qml.dll

qmake
nmake release

:: ----------------------------------------------------------------------------
:: MinGW: generate libleveldb.a

dlltool -dllname go-qml.dll --def capi.def --output-lib libgo-qml.dll.a
copy go-qml.dll ..

:: ----------------------------------------------------------------------------
:: PAUSE


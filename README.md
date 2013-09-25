# QML support for the Go language

This is an ALPHA release
------------------------

This package is in an alpha stage, and still in heavy development. APIs
may change, and things may break.

At this time contributors and developers that are interested in tracking
the development closely are encouraged to use it. If you'd prefer a more
stable release, please hold on a bit and subscribe to the mailing list
for news. It's in a pretty good state, so it shall not take too long.


Demo
----

See [this video](https://www.youtube.com/watch?v=FVQlMrPa7lI) for a quick introduction.


Community
---------

Please join the [mailing list](https://groups.google.com/forum/#!forum/go-qml) for
following relevant development news and discussing project details.


API documentation
------------------

The [API documentation](http://godoc.org/github.com/niemeyer/qml) is available in the usual location.


Installation
------------

To try the alpha release you'll need:

  * Go 1.2 (release candidate), for the C++ support of _go build_
  * Qt 5.0.X or 5.1.X with the development files
  * The Qt headers qmetaobject_p.h and qmetaobjectbuilder_p.h, for the dynamic meta object support

See below for more details about getting these requirements installed in different environments and operating systems.

After the requirements are satisfied, _go get_ should work as usual:

    go get github.com/niemeyer/qml


Requirements on Ubuntu
----------------------

If you are using Ubuntu, the [Ubuntu SDK](http://developer.ubuntu.com/get-started/) will take care of the Qt dependencies:

    $ sudo add-apt-repository ppa:ubuntu-sdk-team/ppa
    $ sudo apt-get update
    $ sudo apt-get install ubuntu-sdk qtbase5-private-dev qtdeclarative5-private-dev

and Go 1.2 may be installed using [godeb](http://blog.labix.org/2013/06/15/in-flight-deb-packages-of-go):

    $ # Pick the right one for your system: 386 or amd64
    $ ARCH=amd64
    $ wget -q https://godeb.s3.amazonaws.com/godeb-$ARCH.tar.gz
    $ tar xzvf godeb-$ARCH.tar.gz
    godeb
    $ sudo mv godeb /usr/local/bin
    $ godeb list | head -1
    1.2rc1
    $ godeb install 1.2rc1
    $ go get github.com/niemeyer/qml


Requirements on Mac OS
----------------------

On Mac OS you'll need gcc (not a symlinked clang, as it complains about `-std=c++11`), and
must specify the `CXX`, `PKG_CONFIG_PATH`, and `CGO_CPPFLAGS` environment variables.

Something along these lines should be effective:

    $ brew tap homebrew/versions
    $ brew install gcc48 qt5

    $ export PKG_CONFIG_PATH=`brew --prefix qt5`/lib/pkgconfig
    $ QT5VERSION=`pkg-config --cflags Qt5Core | sed 's/^.*\(5\..\..\).*/\1/g'`
    $ # For "private/qmetaobject_p.h" inclusion
    $ export CGO_CPPFLAGS=-I`brew --prefix qt5`/include/QtCore/$QT5VERSION/QtCore
    $ CXX=g++-4.8 go get github.com/niemeyer/qml


Requirements on Windows
-----------------------

On windows you need to install the following:

* Mingw-gcc 4.8.1
Download the mingw-get-setup.exe from http://www.mingw.org, download and install the mingw32-gcc compiler from within the
setup GUI.

* Qt 5.1.1 or later
Download Qt 5 binary setup installer for Windows (for Mingw 4.8) from http://qt-project.org and install

* Go 1.2rc1 or later
Download binary MSI installer for Windows from http://golang.org and install.

If you fancy building everything yourselves instead of downloading binaries that is entirely possible too :)

Set the following environment variables:

* CPATH         Add Qt include path and the path of the subfolder 'QtCore\5.X.X\QtCore' under the include folder, replaced with your Qt version.
* LIBRARY_PATH  Add Qt lib path
* PATH          And Qt binary path

Assuming you installed Qt in c:\qt\Qt5.1.1\
that would yield adding the following to the env vars (in environment variables dialog):

CPATH += c:\qt\Qt5.1.1\5.1.1\mingw48_32\include;c:\qt\Qt5.1.1\5.1.1\mingw48_32\include\QtCore\5.1.1\QtCore
LIBRARY_PATH += c:\qt\Qt5.1.1\5.1.1\mingw48_32\lib
PATH += c:\qt\Qt5.1.1\5.1.1\mingw48_32\bin

And finally from the command line (you have to reopen the shell for the env var changes to take effect):

$ go get github.com/niemeyer/qml


Requirements everywhere else
----------------------------

If your operating system does not offer these dependencies readily,
you may still have success installing [Go 1.2rc1](https://code.google.com/p/go/downloads/list?can=1&q=go1.2rc1)
and [Qt 5.0.2](http://download.qt-project.org/archive/qt/5.0/5.0.2/)
directly from the upstreams.

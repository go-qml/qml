# QML support for the Go language

This is an ALPHA release
------------------------

This package is in an alpha stage, and still in heavy development. APIs
may change, and things may break.

At this time contributors and developers that are interested in tracking
the development closely are encouraged to use it. If you'd prefer a more
stable release, please hold on a bit and subscribe to the mailing list
for news. It's in a pretty good state, so it shall not take too long.


Demos
-----

These introductory videos demonstrate the use of Go QML:

  * [Initial demo and overview](http://youtu.be/FVQlMrPa7lI)
  * [Initial demo running on an Ubuntu Touch phone](http://youtu.be/HB-3o8Cysec)
  * [Spinning Gopher with Go + QML + OpenGL](http://youtu.be/qkH7_dtOyPk)
  * [SameGame QML tutorial in Go](http://youtu.be/z8noX48hiMI)


Community
---------

Please join the [mailing list](https://groups.google.com/forum/#!forum/go-qml) for
following relevant development news and discussing project details.


Documentation
-------------

The introductory documentation as well as the detailed API documentation is
available at [gopkg.in/qml.v0](http://godoc.org/gopkg.in/qml.v0).


Installation
------------

To try the alpha release you'll need:

  * Go 1.2, for the C++ support of _go build_
  * Qt 5.0.X or 5.1.X with the development files
  * The Qt headers qmetaobject_p.h and qmetaobjectbuilder_p.h, for the dynamic meta object support

See below for more details about getting these requirements installed in different environments and operating systems.

After the requirements are satisfied, _go get_ should work as usual:

    go get gopkg.in/qml.v0


Requirements on Ubuntu
----------------------

If you are using Ubuntu, the [Ubuntu SDK](http://developer.ubuntu.com/get-started/) will take care of the Qt dependencies:

    $ sudo add-apt-repository ppa:ubuntu-sdk-team/ppa
    $ sudo apt-get update
    $ sudo apt-get install ubuntu-sdk qtbase5-private-dev qtdeclarative5-private-dev libqt5opengl5-dev

and Go 1.2 may be installed using [godeb](http://blog.labix.org/2013/06/15/in-flight-deb-packages-of-go):

    $ # Pick the right one for your system: 386 or amd64
    $ ARCH=amd64
    $ wget -q https://godeb.s3.amazonaws.com/godeb-$ARCH.tar.gz
    $ tar xzvf godeb-$ARCH.tar.gz
    godeb
    $ sudo mv godeb /usr/local/bin
    $ godeb install 1.2
    $ go get gopkg.in/qml.v0


Requirements on Mac OS X
------------------------

On Mac OS X you'll need gcc (not a symlinked clang, as it complains about `-std=c++11`), and
must specify the `CXX`, `PKG_CONFIG_PATH`, and `CGO_CPPFLAGS` environment variables.

Something along these lines should be effective:

    $ brew tap homebrew/versions
    $ brew install gcc48 qt5

    $ export PKG_CONFIG_PATH=`brew --prefix qt5`/lib/pkgconfig
    $ QT5VERSION=`pkg-config --modversion Qt5Core`
    $ # For "private/qmetaobject_p.h" inclusion
    $ export CGO_CPPFLAGS=-I`brew --prefix qt5`/include/QtCore/$QT5VERSION/QtCore
    $ CXX=g++-4.8 go get gopkg.in/qml.v0

For Mac OS X Mavericks you may need to use `brew install qt5 --HEAD` and check that QT5VERSION
is something reasonable like `5.2.0`, `ls /usr/local/Cellar/qt5/HEAD/include/QtCore/ | grep '^5'`
should also work.

Requirements on Windows
-----------------------

On Windows you'll need the following:

  * [MinGW gcc](http://sourceforge.net/projects/mingw/files/latest/download) 4.8.1 (install mingw-get and install the gcc from within the setup GUI)
  * [Qt 5.1.1](http://download.qt-project.org/official_releases/qt/5.1/5.1.1/qt-windows-opensource-5.1.1-mingw48_opengl-x86-offline.exe) for MinGW 4.8
  * [Go 1.2rc1](https://code.google.com/p/go/downloads/list?can=1&q=go1.2rc1) for Windows

Then, assuming Qt was installed under `C:\Qt5.1.1\`, set up the following environment variables in the respective configuration:

    CPATH += C:\Qt5.1.1\5.1.1\mingw48_32\include;C:\Qt5.1.1\5.1.1\mingw48_32\include\QtCore\5.1.1\QtCore
    LIBRARY_PATH += C:\Qt5.1.1\5.1.1\mingw48_32\lib
    PATH += C:\Qt5.1.1\5.1.1\mingw48_32\bin

After reopening the shell for the environment changes to take effect, this should work:

    go get gopkg.in/qml.v0


Requirements everywhere else
----------------------------

If your operating system does not offer these dependencies readily,
you may still have success installing [Go 1.2rc1](https://code.google.com/p/go/downloads/list?can=1&q=go1.2rc1)
and [Qt 5.0.2](http://download.qt-project.org/archive/qt/5.0/5.0.2/)
directly from the upstreams.

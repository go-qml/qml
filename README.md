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
available at [gopkg.in/qml.v1](http://godoc.org/gopkg.in/qml.v1).


Installation
------------

To try the alpha release you'll need:

  * Go 1.2, for the C++ support of _go build_
  * Qt 5.0.X or 5.1.X with the development files
  * The Qt headers qmetaobject_p.h and qmetaobjectbuilder_p.h, for the dynamic meta object support

See below for more details about getting these requirements installed in different environments and operating systems.

After the requirements are satisfied, _go get_ should work as usual:

    go get gopkg.in/qml.v1


Requirements on Ubuntu
----------------------

If you are using Ubuntu, the [Ubuntu SDK](http://developer.ubuntu.com/get-started/) will take care of the Qt dependencies:

    $ sudo add-apt-repository ppa:ubuntu-sdk-team/ppa
    $ sudo apt-get update
    $ sudo apt-get install qtdeclarative5-dev qtbase5-private-dev qtdeclarative5-private-dev libqt5opengl5-dev qtdeclarative5-qtquick2-plugin

and Go 1.2 may be installed using [godeb](http://blog.labix.org/2013/06/15/in-flight-deb-packages-of-go):

    $ # Pick the right one for your system: 386 or amd64
    $ ARCH=amd64
    $ wget -q https://godeb.s3.amazonaws.com/godeb-$ARCH.tar.gz
    $ tar xzvf godeb-$ARCH.tar.gz
    godeb
    $ sudo mv godeb /usr/local/bin
    $ godeb install 1.2
    $ go get gopkg.in/qml.v1


Requirements on Ubuntu Touch
----------------------------

After following the [installation instructions](https://wiki.ubuntu.com/Touch/Install) for Ubuntu Touch,
run the following commands to get a working build environment inside the device:

    $ adb shell
    # cd /tmp
    # wget https://github.com/go-qml/qml/raw/master/cmd/ubuntu-touch/setup.sh
    # /bin/bash setup.sh
    # su - phablet
    $

At the end of setup.sh, the phablet user will have GOPATH=$HOME in the environment,
the qml package will be built, and the particle example will be built and run. For
stopping it from the command line, run as the phablet user:

    $ upstart-app-stop gopkg.in.qml.particle-example

for running it again:

    $ upstart-app-launch gopkg.in.qml.particle-example

These commands depend on the following file, installed by setup.sh:

    ~/.local/share/applications/gopkg.in.qml.particle-example.desktop


Requirements on Mac OS X
------------------------

On Mac OS X you'll need QT5. It's easiest to install with Homebrew, a
third-party package management system for OS X.

Installation instructions for Homebrew are here:

    http://brew.sh/

Then, install the qt5 and pkg-config packages:

    $ brew install qt5 pkg-config

Then, force brew to "link" qt5 (this makes it available under /usr/local):

    $ brew link --force qt5

And finally, fetch and install go-qml:

    $ go get gopkg.in/qml.v1


Requirements on Windows
-----------------------

On Windows you'll need the following:

  * [MinGW gcc](http://tdm-gcc.tdragon.net/download)
  * [Qt 5 with OpenGL](http://download.qt-project.org/archive/qt/5.2/5.2.1/) for MSVC2012
  * [Go 1.3](http://golang.org/dl/) for Windows

Then, assuming Qt was installed under `C:\Qt5.2.1\`, set up the following environment variables in the respective configuration:

	set QTDIR=C:\Qt\Qt5.2.1\5.2.1\msvc2012_64_opengl
	set PATH=%QTDIR%\bin;%PATH%

After reopening the msvc shell for the environment changes to take effect, this should work:

	# download qml.v1
	go get -d gopkg.in/qml.v1

	# Go1.3: build shared library
	cd gopkg.in/qml.v1/cpp && build_msvc.bat
	cd gopkg.in/qml.v1/gl/2.0 && build_msvc.bat

	# Go1.4: build shared library
	go generate gopkg.in/qml.v1
	go generate gopkg.in/qml.v1/gl/2.0

	# install pacakge
	go install gopkg.in/qml.v1
	go install gopkg.in/qml.v1/gl/2.0

Try to run `exmaples/*`.

**Screenshot(Qt5.2/MSVC2012/64bit)**

[![](https://raw.githubusercontent.com/chai2010/qml/v1/screenshot/windows/particle.png)](https://github.com/chai2010/qml/blob/v1/examples/particle/main.go)


Requirements everywhere else
----------------------------

If your operating system does not offer these dependencies readily,
you may still have success installing [Go 1.2rc1](https://code.google.com/p/go/downloads/list?can=1&q=go1.2rc1)
and [Qt 5.0.2](http://download.qt-project.org/archive/qt/5.0/5.0.2/)
directly from the upstreams.  Note that you'll likely have to adapt
environment variables to reflect the custom installation path for
these libraries. See the instructions above for examples.

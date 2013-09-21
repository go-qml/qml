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


Requirements
------------

To try the _alpha release_, you'll need:

  * Go 1.2 (in development), for the C++ support of _go build_

  * For Ubuntu:
    * The current [Ubuntu SDK](http://developer.ubuntu.com/get-started/), or equivalent Qt libraries
    * Packages _qtbase5-private-dev_ and _qtdeclarative5-private-dev_ or equivalent header files, for the dynamic meta object support

  * For MacOS (homebrew):
    * gcc (not a symlinked clang, will complain about `unrecognized command line option "-std=c++11"` otherwise)
    * qt5

In practice, this should work for the Qt dependencies:

Ubuntu:

    sudo add-apt-repository ppa:ubuntu-sdk-team/ppa
    sudo apt-get update
    sudo apt-get install ubuntu-sdk qtbase5-private-dev qtdeclarative5-private-dev

MacOS:

    brew tap homebrew/versions
    brew install gcc48 qt5

For installing Go 1.2 from sources, check the [Go documentation](http://golang.org/doc/install/source).

Installation
------------

Once the requirements above are satisfied, _go get_ should work as usual:

    go get github.com/niemeyer/qml

On MacOS you must specify the `CXX`, `PKG_CONFIG_PATH`, and `CGO_CPPFLAGS` environment variables, for example:

    export PKG_CONFIG_PATH=`brew --prefix qt5`/lib/pkgconfig
    QT5VERSION=`pkg-config --cflags Qt5Core | sed 's/^.*\(5\..\..\).*/\1/g'`
    # For "private/qmetaobject_p.h" inclusion
    export CGO_CPPFLAGS=-I`brew --prefix qt5`/include/QtCore/$QT5VERSION/QtCore
    CXX=g++-4.8 go get github.com/niemeyer/qml

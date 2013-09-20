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

  * Go 1.2 (release candidate), for the C++ support of _go build_
  * The current [Ubuntu SDK](http://developer.ubuntu.com/get-started/), or equivalent Qt libraries
  * Packages _qtbase5-private-dev_ and _qtdeclarative5-private-dev_ or equivalent header files, for the dynamic meta object support

In practice, if you are in Ubuntu, this should work for the Qt dependencies:

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

If you're not in Ubuntu and your operating system does not offer these dependencies,
you may have success installing [Go 1.2rc1](https://code.google.com/p/go/downloads/list?can=1&q=go1.2rc1)
and [Qt 5.0.2](http://download.qt-project.org/archive/qt/5.0/5.0.2/) directly from the upstreams.

Installation
------------

Once the requirements above are satisfied, _go get_ should work as usual:

    go get github.com/niemeyer/qml

# Reproducer for go1.6 port problem

SjB's go 1.6 port worked very well for me for pretty much all of my programs, but
I did discover one corner case that causes problems. The TL;DR is that passing a
go type to qml and back to go does not work.

## What I expect:

Basically the old behaviour. After cloning this repo, try this:

```
GODEBUG=cgocheck=0 go run main.go
```

You should see it print, "Successfully called UseGoType()" and display a white rectangle.

## It doesn't quite work so well with the fix

If we use the 1.6 port (I have copied commit 0309d2df1d6572e107b2bd0712da5b517c4a49be here
for your convenience) then it doesn't work quite like it used to:

```
mv vendor_cjb vendor
go run main.go
```

You should see something like:

```
panic: cannot find fold go reference

goroutine 1 [running, locked to thread]:
panic(0x5716e0, 0xc820090070)
	/usr/local/go/src/runtime/panic.go:464 +0x3e6
.../vendor/gopkg.in/qml%2ev1.getFoldFromGoRef(0x7ffc3a3bc044, 0x8aec00)
.../vendor/gopkg.in/qml.v1/bridge.go:230 +0x9e
... cut ...
```

package qml

// #include <stdlib.h>
//
// #include "capi.h"
//
import "C"

import (
	"errors"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"
)

// Engine provides an environment for instantiating QML components.
type Engine struct {
	Common
	values    map[interface{}]*valueFold
	destroyed bool

	savedAddr unsafe.Pointer // addr might be cleared when destroyed, use this after destroyed

	imageProviders map[string]*func(imageId string, width, height int) image.Image
}

var engines = make(map[unsafe.Pointer]*Engine)

// NewEngine returns a new QML engine.
//
// The Destory method must be called to finalize the engine and
// release any resources used.
func NewEngine() *Engine {
	engine := &Engine{values: make(map[interface{}]*valueFold)}
	RunMain(func() {
		engine.engine = engine
		engine.setAddr(C.newEngine(nil))
		engine.savedAddr = engine.addr
		engine.imageProviders = make(map[string]*func(imageId string, width, height int) image.Image)
		engines[engine.addr] = engine
		stats.enginesAlive(+1)
	})
	return engine
}

func (e *Engine) assertValid() {
	if e.destroyed {
		panic("engine already destroyed")
	}
}

// Destroy finalizes the engine and releases any resources used.
// The engine must not be used after calling this method.
//
// It is safe to call Destroy more than once.
func (e *Engine) Destroy() {
	if !e.destroyed {
		RunMain(func() {
			if !e.destroyed {
				e.destroyed = true
				C.delObjectLater(e.addr)
				if len(e.values) == 0 {
					delete(engines, e.addr)
				} else {
					// The engine reference keeps those values alive.
					// The last value destroyed will clear it.
				}
				stats.enginesAlive(-1)
			}
		})
	}
}

// Load loads a new component with the provided location and with the
// content read from r. The location informs the resource name for
// logged messages, and its path is used to locate any other resources
// referenced by the QML content.
//
// Once a component is loaded, component instances may be created from
// the resulting object via its Create and CreateWindow methods.
func (e *Engine) Load(location string, r io.Reader) (Object, error) {
	var cdata *C.char
	var cdatalen C.int

	qrc := strings.HasPrefix(location, "qrc:")
	if qrc {
		if r != nil {
			return nil, fmt.Errorf("cannot load qrc resource while providing data: %s", location)
		}
	} else {
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		if colon, slash := strings.Index(location, ":"), strings.Index(location, "/"); colon == -1 || slash <= colon {
			if filepath.IsAbs(location) {
				location = "file:///" + filepath.ToSlash(location)
			} else {
				dir, err := os.Getwd()
				if err != nil {
					return nil, fmt.Errorf("cannot obtain absolute path: %v", err)
				}
				location = "file:///" + filepath.ToSlash(filepath.Join(dir, location))
			}
		}

		// Workaround issue #84 (QTBUG-41193) by not refering to an existent file.
		if s := strings.TrimPrefix(location, "file:///"); s != location {
			if _, err := os.Stat(filepath.FromSlash(s)); err == nil {
				location = location + "."
			}
		}

		cdata, cdatalen = unsafeBytesData(data)
	}

	var err error
	cloc, cloclen := unsafeStringData(location)
	comp := &Common{engine: e}
	RunMain(func() {
		// TODO The component's parent should probably be the engine.
		comp.setAddr(C.newComponent(e.addr, nilPtr))
		if qrc {
			C.componentLoadURL(comp.addr, cloc, cloclen)
		} else {
			C.componentSetData(comp.addr, cdata, cdatalen, cloc, cloclen)
		}
		message := C.componentErrorString(comp.addr)
		if message != nilCharPtr {
			err = errors.New(strings.TrimRight(C.GoString(message), "\n"))
			C.free(unsafe.Pointer(message))
		}
	})
	if err != nil {
		return nil, err
	}
	return comp, nil
}

// LoadFile loads a component from the provided QML file.
// Resources referenced by the QML content will be resolved relative to its path.
//
// Once a component is loaded, component instances may be created from
// the resulting object via its Create and CreateWindow methods.
func (e *Engine) LoadFile(path string) (Object, error) {
	if strings.HasPrefix(path, "qrc:") {
		return e.Load(path, nil)
	}
	// TODO Test this.
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return e.Load(path, f)
}

// LoadString loads a component from the provided QML string.
// The location informs the resource name for logged messages, and its
// path is used to locate any other resources referenced by the QML content.
//
// Once a component is loaded, component instances may be created from
// the resulting object via its Create and CreateWindow methods.
func (e *Engine) LoadString(location, qml string) (Object, error) {
	return e.Load(location, strings.NewReader(qml))
}

// Context returns the engine's root context.
func (e *Engine) Context() *Context {
	e.assertValid()
	var ctx Context
	ctx.engine = e
	RunMain(func() {
		ctx.setAddr(C.engineRootContext(e.addr))
	})
	return &ctx
}

func (e *Engine) ClearImportPaths() {
	RunMain(func() {
		C.engineClearImportPaths(e.addr)
	})
}

func (e *Engine) AddImportPath(path string) {
	cpath, cpathLen := unsafeStringData(path)
	RunMain(func() {
		C.engineAddImportPath(e.addr, cpath, cpathLen)
	})
}

func (e *Engine) ClearPluginPaths() {
	RunMain(func() {
		C.engineClearPluginPaths(e.addr)
	})
}

func (e *Engine) AddPluginPath(path string) {
	cpath, cpathLen := unsafeStringData(path)
	RunMain(func() {
		C.engineAddPluginPath(e.addr, cpath, cpathLen)
	})
}

func (e *Engine) ClearComponentCache() {
	RunMain(func() {
		C.engineClearComponentCache(e.addr)
	})
}

// TODO ObjectOf is probably still worth it, but turned out unnecessary
//      for GL functionality. Test it properly before introducing it.

// ObjectOf returns the QML Object representation of the provided Go value
// within the e engine.
//func (e *Engine) ObjectOf(value interface{}) Object {
//	// TODO Would be good to preserve identity on the Go side. See unpackDataValue as well.
//	return &Common{
//		engine: e,
//		addr:   wrapGoValue(e, value, cppOwner),
//	}
//}

// AddImageProvider registers f to be called when an image is requested by QML code
// with the specified provider identifier. It is a runtime error to register the same
// provider identifier multiple times.
//
// The imgId provided to f is the requested image source, with the "image:" scheme
// and provider identifier removed. For example, with an image image source of
// "image://myprovider/icons/home.ext", the respective imgId would be "icons/home.ext".
//
// If either the width or the height parameters provided to f are zero, no specific
// size for the image was requested. If non-zero, the returned image should have the
// the provided size, and will be resized if the returned image has a different size.
//
// See the documentation for more details on image providers:
//
//   http://qt-project.org/doc/qt-5.0/qtquick/qquickimageprovider.html
//
func (e *Engine) AddImageProvider(prvId string, f func(imgId string, width, height int) image.Image) {
	if _, ok := e.imageProviders[prvId]; ok {
		panic(fmt.Sprintf("engine already has an image provider with id %q", prvId))
	}
	e.imageProviders[prvId] = &f
	cprvId, cprvIdLen := unsafeStringData(prvId)
	RunMain(func() {
		qprvId := C.newString(cprvId, cprvIdLen)
		defer C.delString(qprvId)
		C.engineAddImageProvider(e.addr, qprvId, unsafe.Pointer(&f))
	})
}

//export hookRequestImage
func hookRequestImage(imageFunc unsafe.Pointer, cid *C.char, cidLen, cwidth, cheight C.int) unsafe.Pointer {
	f := *(*func(imgId string, width, height int) image.Image)(imageFunc)

	id := unsafeString(cid, cidLen)
	width := int(cwidth)
	height := int(cheight)

	img := f(id, width, height)

	var cimage unsafe.Pointer

	rect := img.Bounds()
	width = rect.Max.X - rect.Min.X
	height = rect.Max.Y - rect.Min.Y
	cimage = C.newImage(C.int(width), C.int(height))

	var cbits []byte
	cbitsh := (*reflect.SliceHeader)((unsafe.Pointer)(&cbits))
	cbitsh.Data = (uintptr)((unsafe.Pointer)(C.imageBits(cimage)))
	cbitsh.Len = width * height * 4 // RGBA
	cbitsh.Cap = cbitsh.Len

	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			*(*uint32)(unsafe.Pointer(&cbits[i])) = (a>>8)<<24 | (r>>8)<<16 | (g>>8)<<8 | (b >> 8)
			i += 4
		}
	}
	return cimage
}

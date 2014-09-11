package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	flagRevert = flag.Bool("revert", false, "revert all changes")
)

var (
	// void *gl1_0_funcs();
	// void gl1_0_glViewport(void *_glfuncs, GLint x, GLint y, GLsizei width, GLsizei height);
	// ...
	// void *gles2_funcs();
	// void gles2_glActiveTexture(void *_glfuncs, GLenum texture);
	// ...
	reGLFunc = regexp.MustCompile(`gl(es)?\d+_(\d+_)?[0-9a-zA-Z_]+\(`)
)

func main() {
	flag.Parse()

	matches, err := filepath.Glob(`../*/funcs.h`)
	if err != nil {
		log.Fatal("filepath.Glob: ", err)
	}
	for i := 0; i < len(matches); i++ {
		dirName := matches[i][:len(matches[i])-len("/funcs.h")]
		processFuncsCpp(dirName)
		processGlGo(dirName)
		generatePro(dirName)
		generateDef(dirName)
		generateBat(dirName)
		supportGenCmd(dirName)
		fmt.Printf("%s ok\n", matches[i])
	}
	fmt.Printf("Done\n")
}

func processFuncsCpp(dirName string) {
	data, err := ioutil.ReadFile(dirName + "/funcs.cpp")
	if err != nil {
		log.Fatal("ioutil.ReadFile: ", err)
	}

	if *flagRevert {
		data = bytes.Replace(data, []byte(`// +build !windows`+"\n"), []byte(""), -1)
		err = ioutil.WriteFile(dirName+"/funcs.cpp", data, 0666)
		if err != nil {
			log.Fatal("ioutil.WriteFile: ", err)
		}
		return
	}

	if !strings.Contains(string(data), `// +build !windows`) {
		data = append([]byte(`// +build !windows`+"\n"), data...)
		err = ioutil.WriteFile(dirName+"/funcs.cpp", data, 0666)
		if err != nil {
			log.Fatal("ioutil.WriteFile: ", err)
		}
	}
}

func processGlGo(dirName string) {
	data, err := ioutil.ReadFile(dirName + "/gl.go")
	if err != nil {
		log.Fatal("ioutil.ReadFile: ", err)
	}
	data, _ = format.Source(data)

	oldStr := "// #cgo pkg-config: Qt5Core Qt5OpenGL\n"
	newStr := "// #cgo !windows pkg-config: Qt5Core Qt5OpenGL\n// #cgo windows LDFLAGS: -L./goqgl -lgoqgl_{{.LibSuffix}}\n"

	newStr = strings.Replace(newStr, "{{.LibSuffix}}", libSuffix(dirName), -1)

	if *flagRevert {
		data = bytes.Replace(data, []byte(newStr), []byte(oldStr), -1)
		data, _ = format.Source(data)

		err = ioutil.WriteFile(dirName+"/gl.go", data, 0666)
		if err != nil {
			log.Fatal("ioutil.WriteFile: ", err)
		}
		return
	} else {
		data = bytes.Replace(data, []byte(oldStr), []byte(newStr), -1)
		data, _ = format.Source(data)

		err = ioutil.WriteFile(dirName+"/gl.go", data, 0666)
		if err != nil {
			log.Fatal("ioutil.WriteFile: ", err)
		}
		return
	}
}

func generatePro(dirName string) {
	var pro = `
TEMPLATE = lib
CONFIG  += dll release
CONFIG  -= embed_manifest_exe embed_manifest_dll
QT      += opengl gui
TARGET   = goqgl_{{.LibSuffix}}

DESTDIR = $${PWD}
INCLUDEPATH += ..

HEADERS += ../funcs.h
SOURCES += ../funcs.cpp

DEF_FILE+= ./goqgl.def
`
	pro = strings.Replace(pro, "{{.LibSuffix}}", libSuffix(dirName), -1)

	if *flagRevert {
		os.RemoveAll(dirName + "/goqgl")
		return
	}

	os.MkdirAll(dirName+"/goqgl", 0666)
	err := ioutil.WriteFile(dirName+"/goqgl/goqgl.pro", []byte(pro), 0666)
	if err != nil {
		log.Fatal("ioutil.WriteFile: ", err)
	}
}

func generateDef(dirName string) {
	var defHeader = `
LIBRARY goqgl_{{.LibSuffix}}.dll

EXPORTS
`
	defHeader = strings.Replace(defHeader, "{{.LibSuffix}}", libSuffix(dirName), -1)

	if *flagRevert {
		os.RemoveAll(dirName + "/goqgl")
		return
	}

	data, err := ioutil.ReadFile(dirName + "/funcs.h")
	if err != nil {
		log.Fatal("ioutil.ReadFile: ", err)
	}

	var funcs []string
	for _, line := range strings.Split(string(data), "\n") {
		if s := reGLFunc.FindString(line); s != "" {
			funcs = append(funcs, s[:len(s)-1])
		}
	}
	sort.Strings(funcs)

	var b bytes.Buffer
	fmt.Fprintf(&b, defHeader[1:])
	for _, s := range funcs {
		fmt.Fprintf(&b, "\t%s\n", s)
	}

	os.MkdirAll(dirName+"/goqgl", 0666)
	err = ioutil.WriteFile(dirName+"/goqgl/goqgl.def", b.Bytes(), 0666)
	if err != nil {
		log.Fatal("ioutil.WriteFile: ", err)
	}
}

func generateBat(dirName string) {
	var bat = `
@echo off

cd %~dp0
setlocal

:: NMake: goqgl_{{.LibSuffix}}.dll
qmake
nmake clean
nmake release

:: MinGW: generate libgoqgl_{{.LibSuffix}}.a
dlltool -dllname goqgl_{{.LibSuffix}}.dll --def goqgl.def --output-lib libgoqgl_{{.LibSuffix}}.a

:: install
copy goqgl_{{.LibSuffix}}.dll %QTDIR%\bin
`
	bat = strings.Replace(bat, "{{.LibSuffix}}", libSuffix(dirName), -1)

	if *flagRevert {
		os.RemoveAll(dirName + "/goqgl")
		return
	}

	os.MkdirAll(dirName+"/goqgl", 0666)
	err := ioutil.WriteFile(dirName+"/goqgl/build_msvc.bat", []byte(bat), 0666)
	if err != nil {
		log.Fatal("ioutil.WriteFile: ", err)
	}
}

func supportGenCmd(dirName string) {
	var gen = `
//go:generate cmd /c call goqgl\build_msvc.bat

package GL
`

	if *flagRevert {
		os.Remove(dirName + "/generate_windows.go")
		return
	}

	err := ioutil.WriteFile(dirName+"/generate_windows.go", []byte(gen), 0666)
	if err != nil {
		log.Fatal("ioutil.WriteFile: ", err)
	}
}

func libSuffix(dirName string) string {
	return strings.Replace(filepath.Base(dirName), ".", "_", -1)
}

package main

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type Define struct {
	Name      string
	Value     string
	Heading   string
	LineBlock int
}

type Func struct {
	Name  string
	Type  string
	Addr  int
	Param []Param
}

type Param struct {
	Name  string
	Type  string
	Addr  int
	Array int
}

type Type struct {
	Name    string
	Type    string
	Comment string
}

type Header struct {
	Define []Define
	Func   []Func
	Type   []Type
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	pkg, err := build.Import("github.com/niemeyer/qml/gl", "", build.FindOnly)
	if err != nil {
		return fmt.Errorf("cannot find qml/gl package path: %v", err)
	}
	data, err := ioutil.ReadFile(filepath.Join(pkg.Dir, "gl.h"))
	if err != nil {
		return fmt.Errorf("cannot read header file: %v", err)
	}
	var header Header
	err = parse(string(data), &header)
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, &header)
	if err != nil {
		return fmt.Errorf("cannot execute template: %v", err)
	}
	out, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("cannot format generated Go code: %v", err)
	}
	os.Stdout.Write(out)
	return nil
}

func typeName(ctypeName string) string {
	if !strings.HasPrefix(ctypeName, "GL") || len(ctypeName) < 3 {
		panic("unexpected C type name: " + ctypeName)
	}
	return string(ctypeName[2]-('a'-'A')) + ctypeName[3:]
}

func constName(defineName string) string {
	if '0' <= defineName[3] && defineName[3] <= '9' {
		return "N" + defineName[3:]
	}
	return defineName[3:]
}

var constLineBlock = -1

func constNewLine(lineBlock int) bool {
	if lineBlock == constLineBlock {
		return false
	}
	constLineBlock = lineBlock
	return true
}

func funcName(cfuncName string) string {
	if !strings.HasPrefix(cfuncName, "gl") || len(cfuncName) < 3 {
		panic("unexpected C function name: " + cfuncName)
	}
	return cfuncName[2:]
}

func funcDocName(cfuncName string) string {
	switch cfuncName[len(cfuncName)-1] {
	case 'i', 'f', 'd', 's':
		switch cfuncName[len(cfuncName)-2] {
		case '1', '2', '3', '4':
			return cfuncName[:len(cfuncName)-2]
		}
	}
	return cfuncName
}

func paramName(cparamName string) string {
	if token.Lookup(cparamName) != token.IDENT {
		return cparamName + "_"
	}
	return cparamName
}

func funcParams(f Func) string {
	var buf []byte
	for i, param := range f.Param {
		if i > 0 {
			buf = append(buf, ", "...)
		}
		buf = append(buf, paramName(param.Name)...)
		buf = append(buf, ' ')
		for j := 0; j < param.Addr; j++ {
			buf = append(buf, '*')
		}
		if param.Array > 0 {
			buf = append(buf, '[')
			buf = append(buf, strconv.Itoa(param.Array)...)
			buf = append(buf, ']')
		}
		buf = append(buf, typeName(param.Type)...)
	}
	return string(buf)
}

func funcCallParams(f Func) string {
	var buf []byte
	for i, param := range f.Param {
		if i > 0 {
			buf = append(buf, ", "...)
		}
		buf = append(buf, "C."...)
		buf = append(buf, param.Type...)
		buf = append(buf, '(')
		buf = append(buf, paramName(param.Name)...)
		buf = append(buf, ')')
	}
	return string(buf)
}

// funcSupported returns whether the given function has wrapping
// properly implemented already.
func funcSupported(f Func) bool {
	if f.Addr > 0 {
		return false
	}
	for _, param := range f.Param {
		if param.Type == "GLvoid" {
			return false
		}
		if param.Addr > 0 {
			return false
		}
		if param.Array > 0 {
			return false
		}
	}
	return true
}

var funcs = template.FuncMap{
	"typeName":       typeName,
	"constName":      constName,
	"constNewLine":   constNewLine,
	"funcName":       funcName,
	"funcDocName":    funcDocName,
	"funcParams":     funcParams,
	"funcCallParams": funcCallParams,
	"funcSupported":  funcSupported,
	"paramName":      paramName,
	"repeat":         strings.Repeat,
}

var tmpl = template.Must(template.New("gl").Funcs(funcs).Parse(`

// ** file automatically generated out of gl.h -- do not edit manually **

package gl

// #define GL_GLEXT_PROTOTYPES
// #include "gl.h"
// #include "glext.h"
//
import "C"

type (
{{range $type := $.Type}}{{if $type.Name | ne "GLvoid"}}	{{$type.Name | typeName}} C.{{$type.Name}}{{if $type.Comment}} /* {{$type.Comment}} */{{end}}
{{end}}{{end}})

const (
{{range $define := $.Define}}{{if $define.LineBlock | constNewLine}}
{{end}}{{if $define.Heading}}	// {{$define.Heading}}
{{end}}	{{$define.Name | constName}} = C.{{$define.Name}}
{{end}})
{{ range $func := $.Func }}{{if $func | funcSupported}}
// See http://www.opengl.org/sdk/docs/man2/xhtml/{{$func.Name | funcDocName}}.xml
func {{$func.Name | funcName}}({{funcParams $func}}) {{if $func.Type}}{{repeat "*" $func.Addr}}{{$func.Type | typeName}} {{end}}{
	{{if $func.Type}}return {{$func.Type | typeName}}({{end}}C.{{$func.Name}}({{funcCallParams $func}}){{if $func.Type}}){{end}}
}
{{end}}{{end}}
`))

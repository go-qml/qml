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
	initFuncNameDocCount(&header)
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

var funcNameDocCount = make(map[string]int)

func initFuncNameDocCount(header *Header) {
	for _, f := range header.Func {
		funcNameDocCount[funcNameDocPrefix(f.Name)]++
	}
}

func funcNameDocPrefix(cfuncName string) string {
	k := len(cfuncName)-1
	if cfuncName[k] == 'v' {
		k--
	}
	switch cfuncName[k] {
	case 'i', 'f', 'd', 's', 'b':
		k--
		if cfuncName[k] == 'u' {
			k--
		}
		switch cfuncName[k] {
		case '1', '2', '3', '4':
			k--
		}
	}
	return cfuncName[:k+1]
}

func funcNameDoc(cfuncName string) string {
	prefix := funcNameDocPrefix(cfuncName)
	if funcNameDocCount[prefix] > 1 {
		return prefix
	}
	return cfuncName
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
			buf = append(buf, "[]"...)
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
	var buf bytes.Buffer
	for i, param := range f.Param {
		if i > 0 {
			buf.WriteString(", ")
		}
		if param.Addr > 0 {
			buf.WriteString("(*")
		}
		buf.WriteString("C.")
		buf.WriteString(param.Type)
		if param.Addr > 0 {
			buf.WriteString(")")
		}
		buf.WriteByte('(')
		if param.Addr > 0 {
			buf.WriteString("unsafe.Pointer(&")
		}
		buf.WriteString(paramName(param.Name))
		if param.Addr > 0 {
			buf.WriteString("[0])")
		}
		buf.WriteByte(')')
	}
	return buf.String()
}

func paramName(cparamName string) string {
	if token.Lookup(cparamName) != token.IDENT {
		return cparamName + "_"
	}
	return cparamName
}

func paramMaxLen(f Func, param Param) string {
	if param.Addr == 0 || len(f.Name) < 3 || f.Name[len(f.Name)-1] != 'v' {
		return ""
	}
	switch f.Name[len(f.Name)-2] {
	case 'i', 'f', 'd', 's':
		switch c := f.Name[len(f.Name)-3]; c {
		case '2', '3', '4':
			return string(c)
		}
	}
	return ""
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
		if param.Array > 0 {
			return false
		}
		if param.Addr > 0 && !strings.HasPrefix(f.Name, "glVertex") && !strings.HasPrefix(f.Name, "glNormal") {
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
	"funcNameDoc":    funcNameDoc,
	"funcParams":     funcParams,
	"funcCallParams": funcCallParams,
	"funcSupported":  funcSupported,
	"paramName":      paramName,
	"paramMaxLen":    paramMaxLen,
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

import (
	"unsafe"
)

type (
{{range $type := $.Type}}{{if $type.Name | ne "GLvoid"}}	{{$type.Name | typeName}} C.{{$type.Name}}{{if $type.Comment}} /* {{$type.Comment}} */{{end}}
{{end}}{{end}})

const (
{{range $define := $.Define}}{{if $define.LineBlock | constNewLine}}
{{end}}{{if $define.Heading}}	// {{$define.Heading}}
{{end}}	{{$define.Name | constName}} = C.{{$define.Name}}
{{end}})
{{ range $func := $.Func }}{{if $func | funcSupported}}
// See http://www.opengl.org/sdk/docs/man2/xhtml/{{$func.Name | funcNameDoc}}.xml
func {{$func.Name | funcName}}({{funcParams $func}}) {{if $func.Type}}{{repeat "*" $func.Addr}}{{$func.Type | typeName}} {{end}}{ {{range $param := $func.Param}}{{with $max := paramMaxLen $func $param}}{{if $max}}if len({{$param.Name | paramName}}) > {{$max}} {
		panic("parameter {{$param.Name | paramName}} has incorrect length")
	}{{end}}
	{{end}}{{end}}{{if $func.Type}}return {{$func.Type | typeName}}({{end}}C.{{$func.Name}}({{funcCallParams $func}}){{if $func.Type}}){{end}}
}
{{end}}{{end}}
`))

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

	GoName string
}

type Func struct {
	Name  string
	Type  string
	Addr  int
	Param []Param

	GoName  string
	GoType  string
	DocName string
}

type Param struct {
	Name  string
	Type  string
	Addr  int
	Array int

	GoName string
	GoType string
}

type Type struct {
	Name    string
	Type    string
	Comment string

	GoName string
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
	prepareHeader(&header)
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

func goTypeName(ctypeName string) string {
	if !strings.HasPrefix(ctypeName, "GL") || len(ctypeName) < 3 {
		panic("unexpected C type name: " + ctypeName)
	}
	return string(ctypeName[2]-('a'-'A')) + ctypeName[3:]
}

func prepareHeader(header *Header) {
	funcNameDocCount := make(map[string]int)

	for fi, f := range header.Func {
		docPrefix := funcNameDocPrefix(f.Name)
		if docPrefix != f.Name {
			funcNameDocCount[docPrefix]++
		}

		if !strings.HasPrefix(f.Name, "gl") || len(f.Name) < 3 {
			panic("unexpected C function name: " + f.Name)
		}
		f.GoName = f.Name[2:]
		if f.Type != "void" {
			f.GoType = goTypeName(f.Type)
		}

		for pi, p := range f.Param {
			switch p.Name {
			case "type", "func", "map":
				p.GoName = "gl" + p.Name
			default:
				if token.Lookup(p.Name) != token.IDENT {
					p.GoName = p.Name + "_"
				} else {
					p.GoName = p.Name
				}
			}
			p.GoType = goTypeName(p.Type)
			f.Param[pi] = p
		}
		header.Func[fi] = f
	}

	for fi, f := range header.Func {
		prefix := funcNameDocPrefix(f.Name)
		if funcNameDocCount[prefix] > 1 {
			f.DocName = prefix
		} else {
			f.DocName = f.Name
		}
		header.Func[fi] = f
	}

	for ti, t := range header.Type {
		t.GoName = goTypeName(t.Name)
		header.Type[ti] = t
	}

	for di, d := range header.Define {
		if !strings.HasPrefix(d.Name, "GL") || len(d.Name) < 3 {
			panic("unexpected C define name: " + d.Name)
		}
		if d.Name[3] >= '0' && d.Name[3] <= '9' {
			d.GoName = "N" + d.Name[3:]
		} else {
			d.GoName = d.Name[3:]
		}
		header.Define[di] = d
	}
}

func funcNameDocPrefix(cfuncName string) string {
	k := len(cfuncName) - 1
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

var constLineBlock = -1

func constNewLine(lineBlock int) bool {
	if lineBlock == constLineBlock {
		return false
	}
	constLineBlock = lineBlock
	return true
}


func funcParams(f Func) string {
	var buf bytes.Buffer
	for i, param := range f.Param {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.GoName)
		buf.WriteByte(' ')
		if param.Type == "GLvoid" && param.Addr > 0 {
			buf.WriteString("interface{}")
			continue
		}
		for j := 0; j < param.Addr; j++ {
			buf.WriteString("[]")
		}
		if param.Array > 0 {
			buf.WriteByte('[')
			buf.WriteString(strconv.Itoa(param.Array))
			buf.WriteByte(']')
		}
		buf.WriteString(param.GoType)
	}
	return buf.String()
}

func funcCallParams(f Func) string {
	var buf bytes.Buffer
	for i, param := range f.Param {
		if i > 0 {
			buf.WriteString(", ")
		}
		if param.Type == "GLvoid" {
			buf.WriteString("unsafe.Pointer(")
			buf.WriteString(param.GoName)
			buf.WriteString("_v.Index(0).Addr().Pointer())")
		} else if param.Addr > 0 {
			buf.WriteString("(*C.")
			buf.WriteString(param.Type)
			buf.WriteString(")(unsafe.Pointer(&")
			buf.WriteString(param.GoName)
			buf.WriteString("[0]))")
		} else {
			buf.WriteString("C.")
			buf.WriteString(param.Type)
			buf.WriteByte('(')
			buf.WriteString(param.GoName)
			buf.WriteByte(')')
		}
	}
	return buf.String()
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
		if param.Array > 0 {
			return false
		}
		if param.Addr > 1 {
			return false
		}
	}
	return true
}

var funcs = template.FuncMap{
	"constNewLine":   constNewLine,
	"funcParams":     funcParams,
	"funcCallParams": funcCallParams,
	"funcSupported":  funcSupported,
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
	"reflect"
	"unsafe"
)

type (
{{range $type := $.Type}}{{if $type.Name | ne "GLvoid"}}	{{$type.GoName}} C.{{$type.Name}}{{if $type.Comment}} /* {{$type.Comment}} */{{end}}
{{end}}{{end}})

const (
{{range $define := $.Define}}{{if $define.LineBlock | constNewLine}}
{{end}}{{if $define.Heading}}	// {{$define.Heading}}
{{end}}	{{$define.GoName}} = C.{{$define.Name}}
{{end}})

{{ range $func := $.Func }}{{if $func | funcSupported}}
// See http://www.opengl.org/sdk/docs/man2/xhtml/{{$func.DocName}}.xml
func {{$func.GoName}}({{funcParams $func}}) {{if $func.GoType}}{{repeat "*" $func.Addr}}{{$func.GoType}} {{end}}{
	{{range $param := $func.Param}}
		{{with $max := paramMaxLen $func $param}}
			{{if $max}}
				if len({{$param.GoName}}) > {{$max}} {
					panic("parameter {{$param.GoName}} has incorrect length")
				}
			{{end}}
		{{end}}
		{{if $param.Type | eq "GLvoid"}}
			{{$param.GoName}}_v := reflect.ValueOf({{$param.GoName}})
			if {{$param.GoName}}_v.Kind() != reflect.Slice {
				panic("parameter {{$param.GoName}} must be a slice")
			}
		{{end}}
	{{end}}
	{{if $func.GoType}}return {{$func.GoType}}({{end}}C.{{$func.Name}}({{funcCallParams $func}}){{if $func.GoType}}){{end}}
}
{{end}}{{end}}
`))

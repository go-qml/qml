package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"go/format"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/xmlpath.v2"
)

type Header struct {
	Class string
	Const []Const
	Func  []Func
	Type  []Type

	FeatureFlags []Const

	GLVersionName  string
	GLVersionLabel string
	GoGLTypes      map[string]bool
}

type Const struct {
	Name      string
	Value     string
	Heading   string
	Comment   string
	LineBlock int
	Disabled  bool

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

	Missing bool
}

type Param struct {
	Name  string
	Type  string
	Addr  int
	Array int
	Const bool

	GoName string
	GoType string
}

type Type struct {
	Name    string
	Type    string
	Comment string

	GoName string
}

type glVersion struct {
	api, number, profile string
}

var glVersions = []glVersion{
	{"gl", "1.0", ""},
	{"gl", "1.1", ""},
	{"gl", "1.2", ""},
	{"gl", "1.3", ""},
	{"gl", "1.4", ""},
	{"gl", "1.5", ""},
	{"gl", "2.0", ""},
	{"gl", "2.1", ""},
	{"gl", "3.0", ""},
	{"gl", "3.1", ""},
	{"gl", "3.2", "core"},
	{"gl", "3.2", "compatibility"},
	{"gl", "3.3", "core"},
	{"gl", "3.3", "compatibility"},
	{"gl", "4.0", "core"},
	{"gl", "4.0", "compatibility"},
	{"gl", "4.1", "core"},
	{"gl", "4.1", "compatibility"},
	{"gl", "4.2", "core"},
	{"gl", "4.2", "compatibility"},
	{"gl", "4.3", "core"},
	{"gl", "4.3", "compatibility"},
	{"gles2", "2.0", ""},
}

func (v glVersion) name() string {
	if v.api == "gles2" {
		return "ES2"
	}
	return v.number
}

func (v glVersion) label() string {
	if v.api == "gles2" {
		return "es2"
	}
	name := strings.Replace(v.number, ".", "_", -1)
	if v.profile == "compatibility" {
		return name + "compat"
	}
	return name + v.profile
}

func (v glVersion) qtheader() string {
	if v.api == "gles2" {
		return "qopenglfunctions.h"
	}
	s := "qopenglfunctions_" + strings.Replace(v.number, ".", "_", -1)
	if v.profile != "" {
		s += "_" + v.profile
	}
	return s + ".h"
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: gengl <qt base include path> <output base path>\n")
		os.Exit(1)
	}
	if err := run(args[0], args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func dirnames(path string) ([]string, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	list, err := dir.Readdirnames(0)
	dir.Close()
	return list, err
}

func parseQtHeader(filename string) (*Header, error) {
	classData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read qt header file: %v", err)
	}
	var header Header
	err = parseQt(string(classData), &header)
	if err != nil {
		return nil, err
	}
	return &header, nil
}

func borrowFuncs(header *Header, filename string) error {
	bheader, err := parseQtHeader(filename)
	if err != nil {
		return err
	}

	seen := make(map[string]bool)
	for _, f := range header.Func {
		seen[f.Name] = true
	}
	for _, f := range bheader.Func {
		if !seen[f.Name] {
			f.Missing = true
			header.Func = append(header.Func, f)
		}
	}
	return nil
}

func run(qtdir, outdir string) error {
	consts, err := parseConsts("gl.xml")
	if err != nil {
		return err
	}

	for _, glVersion := range glVersions {
		header, err := parseQtHeader(filepath.Join(qtdir, "QtGui", glVersion.qtheader()))
		if err != nil {
			return err
		}

		header.GLVersionLabel = glVersion.label()
		header.GLVersionName = glVersion.name()
		header.Const = consts[glVersion]

		if glVersion.api == "gles2" {
			// Qt < 5.3 misses several ES2 entries in QOpenGLFunctions. As a workaround,
			// find the missing entries from the pure ES2 class and use them directly.
			err := borrowFuncs(header, filepath.Join(qtdir, "QtGui", "qopenglfunctions_es2.h"))
			if err != nil {
				return err
			}
		}

		err = prepareHeader(header)
		if err != nil {
			return err
		}

		fileContent := make(map[string][]byte)
		for _, pkgFile := range packageFiles {
			var buf bytes.Buffer
			err = pkgFile.Template.Execute(&buf, header)
			if err != nil {
				return fmt.Errorf("cannot execute template: %v", err)
			}
			data := buf.Bytes()
			if strings.HasSuffix(pkgFile.Name, ".go") {
				newdata, err := format.Source(data)
				if err != nil {
					return fmt.Errorf("\n%s\ncannot format generated Go code: %v\n", data, err)
				}
				data = newdata
			}
			fileContent[pkgFile.Name] = data
		}

		glDir := filepath.Join(outdir, strings.Replace(header.GLVersionLabel, "_", ".", -1))
		err = os.MkdirAll(glDir, 0755)
		if err != nil {
			return fmt.Errorf("cannot make package directory at %s: %v", glDir, err)
		}
		oldNames, err := dirnames(glDir)
		if err != nil {
			return fmt.Errorf("cannot list contents of directory %s: %v", glDir, err)
		}
		for _, oldName := range oldNames {
			oldPath := filepath.Join(glDir, oldName)
			if err := os.Remove(oldPath); err != nil {
				return fmt.Errorf("cannot remove previous file in %s: %v", glDir, err)
			}
		}
		for name, data := range fileContent {
			err = ioutil.WriteFile(filepath.Join(glDir, name), data, 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func parseRegistry() (*xmlpath.Node, error) {
	glxml, err := os.Open("gl.xml")
	if err != nil {
		return nil, fmt.Errorf("cannot open gl.xml file: %v", err)
	}
	defer glxml.Close()

	root, err := xmlpath.Parse(glxml)
	if err != nil {
		return nil, fmt.Errorf("cannot parse gl.xml file: %v", err)
	}
	return root, nil
}

type glRegistry struct {
	Enums    []glEnum    `xml:"enums>enum"`
	Groups   []glGroup   `xml:"groups>group"`
	Features []glFeature `xml:"feature"`
}

type glFeature struct {
	API      string     `xml:"api,attr"`
	Number   string     `xml:"number,attr"`
	Requires []glChange `xml:"require"`
	Removes  []glChange `xml:"remove"`
}

type glChange struct {
	Profile string   `xml:"profile,attr"`
	Enums   []glEnum `xml:"enum"`
}

type glGroup struct {
	Name  string   `xml:"name,attr"`
	Enums []glEnum `xml:"enum"`
}

type glEnum struct {
	API   string `xml:"api,attr"`
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type groupTweak struct {
	group   string
	rename  string
	replace []string
	append  []string
	reorder bool
}

var groupTweaks = []groupTweak{{
	group:   "Boolean",
	reorder: true,
	append:  []string{"GL_NONE"},
}, {
	group:   "DataType",
	reorder: true,
	replace: []string{
		"GL_BYTE",
		"GL_UNSIGNED_BYTE",
		"GL_SHORT",
		"GL_UNSIGNED_SHORT",
		"GL_INT",
		"GL_UNSIGNED_INT",
		"GL_FLOAT",
		"GL_2_BYTES",
		"GL_2_BYTES_NV",
		"GL_3_BYTES",
		"GL_3_BYTES_NV",
		"GL_4_BYTES",
		"GL_4_BYTES_NV",
		"GL_DOUBLE",
		"GL_DOUBLE_EXT",
		"GL_HALF_FLOAT",
		"GL_HALF_FLOAT_ARB",
		"GL_HALF_FLOAT_NV",
		"GL_HALF_APPLE",
		"GL_FIXED",
		"GL_FIXED_OES",
		"GL_INT64_NV",
		"GL_UNSIGNED_INT64_ARB",
		"GL_UNSIGNED_INT64_NV",
	},
}, {
	group: "BlendingFactorSrc",
	replace: []string{
		"GL_CONSTANT_ALPHA_EXT",
		"GL_CONSTANT_COLOR_EXT",
		"GL_DST_ALPHA",
		"GL_DST_COLOR",
		"GL_ONE",
		"GL_ONE_MINUS_DST_ALPHA",
		"GL_ONE_MINUS_DST_COLOR",
		"GL_ONE_MINUS_SRC_ALPHA",
		"GL_ONE_MINUS_SRC_COLOR",
		"GL_SRC_ALPHA",
		"GL_SRC_COLOR",
		"GL_SRC_ALPHA_SATURATE",
		"GL_ZERO",
	},
}, {
	group: "AttribMask",
	append: []string{
		"GL_COVERAGE_BUFFER_BIT_NV",
	},
}}

func tweakGroups(groups []glGroup) {
	tweaks := make(map[string]groupTweak)
	reorder := make(map[string]int)
	top := 0
	for _, tweak := range groupTweaks {
		tweaks[tweak.group] = tweak
		if tweak.reorder {
			reorder[tweak.group] = top
			top++
		}
	}

	// Take reordered groups out, leaving space at the start.
	stash := make([]glGroup, top)
	for i, group := range groups {
		newi, ok := reorder[group.Name]
		if !ok {
			continue
		}
		stash[newi] = group
		copy(groups[1:i+1], groups[0:i])
	}
	// Put reordered groups back, in the proper order.
	copy(groups, stash)

	for i, group := range groups {
		tweak, ok := tweaks[group.Name]
		if !ok {
			continue
		}
		if tweak.rename != "" {
			group.Name = tweak.rename
		}
		if tweak.replace != nil {
			group.Enums = group.Enums[:0]
			for _, name := range tweak.replace {
				group.Enums = append(group.Enums, glEnum{Name: name})
			}
		}
		if tweak.append != nil {
			for _, name := range tweak.append {
				group.Enums = append(group.Enums, glEnum{Name: name})
			}
		}
		groups[i] = group
	}
}

func parseConsts(filename string) (map[glVersion][]Const, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %v", filename, err)
	}
	defer f.Close()

	var registry glRegistry
	err = xml.NewDecoder(f).Decode(&registry)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %s: %v", filename, err)
	}

	var all = make(map[glVersion]map[string]bool)
	var last = make(map[string]map[string]bool)
	for _, feature := range registry.Features {
		profiles := make(map[string]bool)
		for _, require := range feature.Requires {
			profiles[require.Profile] = true
		}
		for _, remove := range feature.Removes {
			profiles[remove.Profile] = true
		}

		for profile := range profiles {
			required := make(map[string]bool)
			inherit, ok := last[feature.API+":"+profile]
			if !ok && profile != "" {
				inherit = last[feature.API+":"]
			}
			for name := range inherit {
				required[name] = true
			}

			for _, require := range feature.Requires {
				if require.Profile == profile || profile == "" {
					for _, enum := range require.Enums {
						required[enum.Name] = true
					}
				}
			}
			for _, remove := range feature.Removes {
				if remove.Profile == profile || profile == "" {
					for _, enum := range remove.Enums {
						delete(required, enum.Name)
					}
				}
			}

			all[glVersion{feature.API, feature.Number, profile}] = required
			last[feature.API+":"+profile] = required
		}
	}

	// Enums in groups and requires/removes have no values.
	enums := make(map[string]glEnum)
	for _, enum := range registry.Enums {
		enums[enum.Name] = enum
	}

	tweakGroups(registry.Groups)

	consts := make(map[glVersion][]Const)
	for _, glVersion := range glVersions {
		var required = all[glVersion]
		var done = make(map[string]bool)
		var lineblock = 0
		var vconsts []Const
		for _, group := range registry.Groups {
			for _, enum := range group.Enums {
				if required[enum.Name] && !done[enum.Name] {
					done[enum.Name] = true
					c := Const{
						Name:      enum.Name,
						Value:     enums[enum.Name].Value,
						LineBlock: lineblock,
					}
					vconsts = append(vconsts, c)
				}
			}
			lineblock++
		}
		// Everything else not found in groups.
		for _, enum := range registry.Enums {
			if required[enum.Name] && !done[enum.Name] {
				done[enum.Name] = true
				c := Const{
					Name:      enum.Name,
					Value:     enum.Value,
					LineBlock: lineblock,
				}
				vconsts = append(vconsts, c)
			}
		}
		consts[glVersion] = vconsts
	}
	// Version 1.0 has no enums. Copy from 1.1.
	consts[glVersion{"gl", "1.0", ""}] = consts[glVersion{"gl", "1.1", ""}]
	return consts, nil
}

func goTypeName(ctypeName string) string {
	if !strings.HasPrefix(ctypeName, "GL") || len(ctypeName) < 3 {
		panic("unexpected C type name: " + ctypeName)
	}
	return string(ctypeName[2]-('a'-'A')) + ctypeName[3:]
}

var enumsPath = xmlpath.MustCompile("/registry/enums/enum")

func prepareHeader(header *Header) error {
	funcNameDocCount := make(map[string]int)

	header.GoGLTypes = make(map[string]bool)

	for fi, f := range header.Func {
		docPrefix := funcNameDocPrefix(f.Name)
		if docPrefix != f.Name {
			funcNameDocCount[docPrefix]++
		}

		if !strings.HasPrefix(f.Name, "gl") || len(f.Name) < 3 {
			panic("unexpected C function name: " + f.Name)
		}
		f.GoName = f.Name[2:]
		if f.Type == "int" {
			// Some consistency. It's in a gl* function after all.
			f.Type = "GLint"
		}
		if f.Type != "void" {
			f.GoType = goTypeName(f.Type)
		}
		header.GoGLTypes[f.Type] = true

		for pi, p := range f.Param {
			switch p.Name {
			case "type", "func", "map", "string":
				p.GoName = "gl" + p.Name
			default:
				if token.Lookup(p.Name) != token.IDENT {
					p.GoName = p.Name + "_"
				} else {
					p.GoName = p.Name
				}
			}
			// Some consistency. Those are a gl* function after all.
			switch p.Type {
			case "void":
				p.Type = "GLvoid"
			case "char":
				p.Type = "GLchar"
			case "qopengl_GLsizeiptr", "qopengl_GLintptr":
				p.Type = p.Type[8:]
			}
			p.GoType = goTypeName(p.Type)
			header.GoGLTypes[p.Type] = true
			f.Param[pi] = p
		}
		header.Func[fi] = f
	}

	delete(header.GoGLTypes, "void")
	delete(header.GoGLTypes, "GLvoid")

	for fi, f := range header.Func {
		prefix := funcNameDocPrefix(f.Name)
		if true || funcNameDocCount[prefix] > 1 {
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

	for di, d := range header.Const {
		if !strings.HasPrefix(d.Name, "GL") || len(d.Name) < 3 {
			panic("unexpected C define name: " + d.Name)
		}
		if d.Name[3] >= '0' && d.Name[3] <= '9' {
			d.GoName = "N" + d.Name[3:]
		} else {
			d.GoName = d.Name[3:]
		}
		header.Const[di] = d
	}

	return nil
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

type funcTweak struct {
	params   string
	preamble string
}

var funcTweaks = map[string]funcTweak{
	"glShaderSource": {
		params:   "shader Uint, source [][]byte",
		preamble: `
			count := len(source)
			length := make([]Int, count)
			glstring := make([]unsafe.Pointer, count)
			for i, buf := range source {
				length[i] = Int(len(buf))
				if len(buf) > 0 {
					glstring[i] = unsafe.Pointer(&buf[0])
				} else {
					glstring[i] = unsafe.Pointer(uintptr(0))
				}
			}
		`,
	},
}

func funcParams(f Func) string {
	if tweak, ok := funcTweaks[f.Name]; ok && tweak.params != "" {
		return tweak.params
	}
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

func funcCallPreamble(f Func) string {
	return funcTweaks[f.Name].preamble
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
			buf.WriteByte('(')
			for i := 0; i < param.Addr; i++ {
				buf.WriteByte('*')
			}
			buf.WriteString("C.")
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

func funcCParams(f Func) string {
	var buf bytes.Buffer
	for i, param := range f.Param {
		if i > 0 {
			buf.WriteString(", ")
		}
		if param.Const {
			buf.WriteString("const ")
		}
		buf.WriteString(param.Type)
		for j := 0; j < param.Addr; j++ {
			buf.WriteString("*")
		}
		if param.Array > 0 {
			buf.WriteByte('[')
			buf.WriteString(strconv.Itoa(param.Array))
			buf.WriteByte(']')
		}
		buf.WriteByte(' ')
		buf.WriteString(param.GoName)
	}
	return buf.String()
}

func funcCCallParams(f Func) string {
	var buf bytes.Buffer
	for i, param := range f.Param {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.GoName)
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
	if _, ok := funcTweaks[f.Name]; ok {
		return true
	}
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
	"constNewLine":     constNewLine,
	"funcParams":       funcParams,
	"funcCallParams":   funcCallParams,
	"funcCallPreamble": funcCallPreamble,
	"funcCParams":      funcCParams,
	"funcCCallParams":  funcCCallParams,
	"funcSupported":    funcSupported,
	"paramMaxLen":      paramMaxLen,
	"repeat":           strings.Repeat,
	"tolower":          strings.ToLower,
	"goTypeName":       goTypeName,
}

type packageFile struct {
	Name     string
	Template *template.Template
}

var packageFiles = []packageFile{
	{"gl.go", tmplGo},
	{"funcs.cpp", tmplFuncsCpp},
	{"funcs.h", tmplFuncsH},
}

var tmplGo = template.Must(template.New("gl").Funcs(funcs).Parse(`
// ** file automatically generated by glgen -- do not edit manually **

package GL

// #cgo CXXFLAGS: -std=c++0x -pedantic-errors -Wall -fno-strict-aliasing 
// #cgo LDFLAGS: -lstdc++
// #cgo !darwin LDFLAGS: -lGL
// #cgo  darwin LDFLAGS: -framework OpenGL
// #cgo pkg-config: Qt5Core Qt5OpenGL
//
// #include "funcs.h"
//
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"

	"gopkg.in/qml.v1/gl"
)

type (
{{range $gltype, $foo := $.GoGLTypes}}{{goTypeName $gltype}} C.{{$gltype}}
{{end}})

// API returns a value that offers methods matching the OpenGL version {{$.GLVersionName}} API.
//
// The returned API must not be used after the provided OpenGL context becomes invalid.
func API(context gl.Contexter) *GL {
	gl := &GL{}
	gl.funcs = C.gl{{$.GLVersionLabel}}_funcs()
	if gl.funcs == nil {
		panic(fmt.Errorf("OpenGL version {{$.GLVersionName}} is not available"))
	}
	return gl
}

// GL implements the OpenGL version {{$.GLVersionName}} API. Values of this
// type must be created via the API function, and it must not be used after
// the associated OpenGL context becomes invalid.
type GL struct {
	funcs unsafe.Pointer
}

const ({{range $const := $.Const}}{{if $const.LineBlock | constNewLine}}
{{end}}{{if $const.Heading}}	// {{$const.Heading}}
{{end}}	{{if $const.Disabled}}//{{end}}{{$const.GoName}} = {{$const.Value}}{{if $const.Comment}}	// {{$const.Comment}}{{end}}
{{end}})

{{ range $func := $.Func }}{{if $func | funcSupported}}
// See http://www.opengl.org/sdk/docs/man3/xhtml/{{$func.DocName}}.xml
func (gl *GL) {{$func.GoName}}({{funcParams $func}}) {{if $func.GoType}}{{repeat "*" $func.Addr}}{{$func.GoType}} {{end}}{ {{/*
*/}}	{{funcCallPreamble $func}} {{/*
*/}}	{{range $param := $func.Param}} {{/*
*/}}		{{with $max := paramMaxLen $func $param}} {{/*
*/}}			{{if $max}} {{/*
*/}}				if len({{$param.GoName}}) > {{$max}} {
					panic("parameter {{$param.GoName}} has incorrect length")
				}
			{{end}}
		{{end}} {{/*
*/}}		{{if $param.Type | eq "GLvoid"}} {{/*
*/}}			{{$param.GoName}}_v := reflect.ValueOf({{$param.GoName}})
			if {{$param.GoName}}_v.Kind() != reflect.Slice {
				panic("parameter {{$param.GoName}} must be a slice")
			}
		{{end}} {{/*
*/}}	{{end}} {{/*
*/}}	{{if $func.GoType}}return {{$func.GoType}}({{end}}C.gl{{$.GLVersionLabel}}_{{$func.Name}}(gl.funcs{{if $func.Param}}, {{funcCallParams $func}}{{end}}){{if $func.GoType}}){{end}}
}
{{end}}{{end}}
`))

var tmplFuncsCpp = template.Must(template.New("gl").Funcs(funcs).Parse(`
// ** file automatically generated by glgen -- do not edit manually **

#include <QOpenGLContext>
#include <QtGui/{{$.Class | tolower}}.h>

#include "funcs.h"

void *gl{{$.GLVersionLabel}}_funcs() {
	{{$.Class}}* funcs = QOpenGLContext::currentContext()->{{if eq $.Class "QOpenGLFunctions"}}functions{{else}}versionFunctions<{{$.Class}}>{{end}}();
	if (!funcs) {
		return 0;
	}{{if ne $.Class "QOpenGLFunctions"}}
	funcs->initializeOpenGLFunctions();{{end}}
	return funcs;
}

{{ range $func := $.Func }}{{if $func | funcSupported}}
{{$func.Type}} gl{{$.GLVersionLabel}}_{{$func.Name}}(void *_glfuncs{{if $func.Param}}, {{funcCParams $func}}{{end}})
{
	{{if not $func.Missing}}{{$.Class}}* _qglfuncs = reinterpret_cast<{{$.Class}}*>(_glfuncs);
	{{end}}{{if $func.GoType}}return {{end}}{{if not $func.Missing}}_qglfuncs->{{end}}{{$func.Name}}({{funcCCallParams $func}});
}
{{end}}{{end}}
`))

var tmplFuncsH = template.Must(template.New("gl").Funcs(funcs).Parse(`
// ** file automatically generated by glgen -- do not edit manually **

#ifndef __cplusplus
#include <inttypes.h>
#include <stddef.h>
typedef unsigned int	GLenum;
typedef unsigned char	GLboolean;
typedef unsigned int	GLbitfield;
typedef void		GLvoid;
typedef char            GLchar;
typedef signed char	GLbyte;		/* 1-byte signed */
typedef short		GLshort;	/* 2-byte signed */
typedef int		GLint;		/* 4-byte signed */
typedef unsigned char	GLubyte;	/* 1-byte unsigned */
typedef unsigned short	GLushort;	/* 2-byte unsigned */
typedef unsigned int	GLuint;		/* 4-byte unsigned */
typedef int		GLsizei;	/* 4-byte signed */
typedef float		GLfloat;	/* single precision float */
typedef float		GLclampf;	/* single precision float in [0,1] */
typedef double		GLdouble;	/* double precision float */
typedef double		GLclampd;	/* double precision float in [0,1] */
typedef int64_t         GLint64;
typedef uint64_t        GLuint64;
typedef ptrdiff_t       GLintptr;
typedef ptrdiff_t       GLsizeiptr;
typedef ptrdiff_t       GLintptrARB;
typedef ptrdiff_t       GLsizeiptrARB;
typedef struct __GLsync *GLsync;
#endif

#ifdef __cplusplus
extern "C" {
#endif

void *gl{{$.GLVersionLabel}}_funcs();

{{ range $func := $.Func }}{{if $func | funcSupported}}{{$func.Type}} gl{{$.GLVersionLabel}}_{{$func.Name}}(void *_glfuncs{{if $func.Param}}, {{funcCParams $func}}{{end}});
{{end}}{{end}}

#ifdef __cplusplus
} // extern "C"
#endif
`))

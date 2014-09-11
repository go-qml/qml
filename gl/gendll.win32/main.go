// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	// void *gl1_0_funcs();
	// void gl1_0_glViewport(void *_glfuncs, GLint x, GLint y, GLsizei width, GLsizei height);
	// ...
	// void *gles2_funcs();
	// void gles2_glActiveTexture(void *_glfuncs, GLenum texture);
	// ...
	re = regexp.MustCompile(`gl(es)?\d+_(\d+_)?[0-9a-zA-Z_]+\(`)
)

func main() {
	matches, err := filepath.Glob(`../*/funcs.h`)
	if err != nil {
		log.Fatal("filepath.Glob: ", err)
	}
	for i := 0; i < len(matches); i++ {
		dirName := matches[i][:len(matches[i])-len("/funcs.h")]
		processFuncsCpp(dirName)
		makeDef(dirName)
		fmt.Printf("%s ok\n", matches[i])
	}
	fmt.Printf("Done\n")
}

func processFuncsCpp(dirName string) {
	data, err := ioutil.ReadFile(dirName + "/funcs.cpp")
	if err != nil {
		log.Fatal("ioutil.ReadFile: ", err)
	}
	if !strings.Contains(string(data), `// +build !windows`) {
		data = append([]byte(`// +build !windows`+"\n"), data...)
		err = ioutil.WriteFile(dirName+"/funcs.cpp", data, 0666)
		if err != nil {
			log.Fatal("ioutil.WriteFile: ", err)
		}
	}
}

func makeDef(dirName string) {
	data, err := ioutil.ReadFile(dirName + "/funcs.h")
	if err != nil {
		log.Fatal("ioutil.ReadFile: ", err)
	}

	var funcs []string
	for _, line := range strings.Split(string(data), "\n") {
		if s := re.FindString(line); s != "" {
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
	os.Remove(dirName + "/goqgl.def")

	err = ioutil.WriteFile(dirName+"/goqgl/goqgl.def", b.Bytes(), 0666)
	if err != nil {
		log.Fatal("ioutil.WriteFile: ", err)
	}
}

var defHeader = `
; Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
; Use of this source code is governed by a BSD-style
; license that can be found in the LICENSE file.

; Auto Genrated by makedef.go; DO NOT EDIT!!

LIBRARY goqgl.dll

EXPORTS
`

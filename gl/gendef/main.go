// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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
		makeDef(matches[i])
		fmt.Printf("%s ok\n", matches[i])
	}
	fmt.Printf("Done\n")
}

func makeDef(filename string) {
	data, err := ioutil.ReadFile(filename)
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
	fmt.Fprintf(&b, header[1:])
	for _, s := range funcs {
		fmt.Fprintf(&b, "\t%s\n", s)
	}

	defName := strings.Replace(filename, "funcs.h", "goqgl.def", -1)
	err = ioutil.WriteFile(defName, b.Bytes(), 0666)
	if err != nil {
		log.Fatal("ioutil.WriteFile: ", err)
	}
}

var header = `
; Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
; Use of this source code is governed by a BSD-style
; license that can be found in the LICENSE file.

; Auto Genrated by makedef.go; DO NOT EDIT!!

LIBRARY goqgl.dll

EXPORTS
`

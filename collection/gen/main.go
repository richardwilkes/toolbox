// +build gen

package main

import (
	"bytes"
	"go/format"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/log/jot"
)

//go:generate go run main.go

const codeGenNotice = "// Code generated - DO NOT EDIT.\n"

func main() {
	data, err := ioutil.ReadFile("set.go.tmpl")
	jot.FatalIfErr(err)
	tmpl := template.Must(template.New("gen").Funcs(template.FuncMap{"lower": strings.ToLower}).Parse(string(data)))
	for _, one := range []string{
		"Byte",
		"Complex64",
		"Complex128",
		"Float32",
		"Float64",
		"Int",
		"Int8",
		"Int16",
		"Int32",
		"Int64",
		"Rune",
		"String",
		"Uint",
		"Uint8",
		"Uint16",
		"Uint32",
		"Uint64",
	} {
		var buffer bytes.Buffer
		buffer.WriteString(codeGenNotice)
		jot.FatalIfErr(tmpl.Execute(&buffer, one))
		d, err := format.Source(buffer.Bytes())
		jot.FatalIfErr(err)
		jot.FatalIfErr(ioutil.WriteFile("../"+strings.ToLower(one)+"set_gen.go", d, 0644))
	}
	atexit.Exit(0)
}

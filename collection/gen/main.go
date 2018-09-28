// +build gen

package main

import (
	"bufio"
	"html/template"
	"io/ioutil"
	"os"
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
		f, err := os.Create("../" + strings.ToLower(one) + "set_gen.go")
		jot.FatalIfErr(err)
		w := bufio.NewWriter(f)
		w.WriteString(codeGenNotice)
		jot.FatalIfErr(tmpl.Execute(w, one))
		jot.FatalIfErr(w.Flush())
		jot.FatalIfErr(f.Close())
	}
	atexit.Exit(0)
}

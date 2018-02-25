// +build gen

package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

//go:generate go run main.go

const codeGenNotice = "// Code generated - DO NOT EDIT.\n"

func main() {
	data, err := ioutil.ReadFile("set.go.tmpl")
	failOnErr(err)
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
		failOnErr(err)
		w := bufio.NewWriter(f)
		w.WriteString(codeGenNotice)
		failOnErr(tmpl.Execute(w, one))
		failOnErr(w.Flush())
		failOnErr(f.Close())
	}
}

func failOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
}

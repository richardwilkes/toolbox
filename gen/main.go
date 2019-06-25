package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/txt"
)

//go:generate go run main.go

type cmdlineInfo struct {
	Type           string
	Parser         string
	NeedConversion bool
}

var (
	setTypes = []string{
		"byte",
		"complex64",
		"complex128",
		"float32",
		"float64",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"rune",
		"string",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
	}
	cmdlineTypes = []cmdlineInfo{
		{"bool", "strconv.ParseBool(str)", false},
		{"int", "strconv.ParseInt(str, 0, 64)", true},
		{"int8", "strconv.ParseInt(str, 0, 8)", true},
		{"int16", "strconv.ParseInt(str, 0, 16)", true},
		{"int32", "strconv.ParseInt(str, 0, 32)", true},
		{"int64", "strconv.ParseInt(str, 0, 64)", false},
		{"uint", "strconv.ParseUint(str, 0, 64)", true},
		{"uint8", "strconv.ParseUint(str, 0, 8)", true},
		{"uint16", "strconv.ParseUint(str, 0, 16)", true},
		{"uint32", "strconv.ParseUint(str, 0, 32)", true},
		{"uint64", "strconv.ParseUint(str, 0, 64)", false},
		{"float32", "strconv.ParseFloat(str, 32)", true},
		{"float64", "strconv.ParseFloat(str, 64)", false},
		{"string", "str, error(nil)", false},
		{"time.Duration", "time.ParseDuration(str)", false},
	}
)

func main() {
	for _, one := range collectGenFiles() {
		jot.FatalIfErr(os.Remove(one))
	}
	tmpl := template.New("").Funcs(template.FuncMap{
		"first_to_upper": txt.FirstToUpper,
		"name":           toName,
	})
	tmpls, err := tmpl.ParseGlob("tmpl/*.go.tmpl")
	jot.FatalIfErr(errs.Wrap(err))
	for _, one := range setTypes {
		jot.FatalIfErr(writeGoTemplate(tmpls, "set.go.tmpl", "../collection/"+one+"set_gen.go", one))
	}
	for _, one := range cmdlineTypes {
		jot.FatalIfErr(writeGoTemplate(tmpls, "values.go.tmpl", "../cmdline/"+toName(one.Type)+"_value_gen.go", one))
	}
	atexit.Exit(0)
}

func toName(in string) string {
	if i := strings.Index(in, "."); i != -1 {
		return strings.ToLower(in[i+1:])
	}
	return in
}

func collectGenFiles() []string {
	var result []string
	rootPath, err := filepath.Abs("..")
	jot.FatalIfErr(err)
	jot.FatalIfErr(filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info != nil {
			name := strings.ToLower(info.Name())
			if info.IsDir() {
				if path != rootPath && (name == ".git" || name == ".cvs") {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasSuffix(name, "_gen.go") {
				result = append(result, path)
			}
		}
		return nil
	}))
	sort.Slice(result, func(i, j int) bool { return txt.NaturalLess(result[i], result[j], true) })
	return result
}

func writeGoTemplate(tmpls *template.Template, tmpl, outputPath string, tmplArg interface{}) error {
	var buffer bytes.Buffer
	// I intentionally do not include the standard "Code generated ... DO NOT EDIT"
	// comment at the top of the file because I want linting to run on the
	// file. Instead, this comment is inserted, which is similar, but doesn't
	// trigger the linter to ignore the file.
	fmt.Fprintf(&buffer, "// Code created from %q - don't edit by hand\n\n", tmpl)
	if err := tmpls.ExecuteTemplate(&buffer, tmpl, tmplArg); err != nil {
		return errs.Wrap(err)
	}
	data := buffer.Bytes()
	if formatted, err := format.Source(data); err != nil {
		fmt.Println(errs.NewWithCause(fmt.Sprintf("unable to format %q", outputPath), err))
	} else {
		data = formatted
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return errs.Wrap(err)
	}
	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

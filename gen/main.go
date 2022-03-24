// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package main

import (
	"bytes"
	"fmt"
	"go/format"
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

type fixedTestInfo struct {
	Bits   int
	Digits int
}

var (
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
	fixed64Digits  = []int{2, 3, 4, 6}
	fixed128Digits = []int{2, 3, 4, 6, 16}
)

func main() {
	for _, one := range collectGenFiles() {
		jot.FatalIfErr(os.Remove(one))
	}
	tmpl := template.New("").Funcs(template.FuncMap{
		"first_to_upper": txt.FirstToUpper,
		"name":           toName,
		"wrap_comment":   wrapCommentWithLength,
		"repeat":         strings.Repeat,
		"add":            add,
		"sub":            sub,
	})
	tmpls, err := tmpl.ParseGlob("tmpl/*.go.tmpl")
	jot.FatalIfErr(errs.Wrap(err))
	for _, one := range cmdlineTypes {
		jot.FatalIfErr(writeGoTemplate(tmpls, "values.go.tmpl", "../cmdline/"+toName(one.Type)+"_value_gen.go", one))
	}
	for _, one := range fixed64Digits {
		jot.FatalIfErr(writeGoTemplate(tmpls, "fixed64.go.tmpl", fmt.Sprintf("../xmath/fixed/F64d%d_gen.go", one), one))
		jot.FatalIfErr(writeGoTemplate(tmpls, "fixed_test.go.tmpl", fmt.Sprintf("../xmath/fixed/F64d%d_gen_test.go", one), &fixedTestInfo{Bits: 64, Digits: one}))
	}
	for _, one := range fixed128Digits {
		jot.FatalIfErr(writeGoTemplate(tmpls, "fixed128.go.tmpl", fmt.Sprintf("../xmath/fixed/F128d%d_gen.go", one), one))
		jot.FatalIfErr(writeGoTemplate(tmpls, "fixed_test.go.tmpl", fmt.Sprintf("../xmath/fixed/F128d%d_gen_test.go", one), &fixedTestInfo{Bits: 128, Digits: one}))
	}
	atexit.Exit(0)
}

func toName(in string) string {
	if i := strings.Index(in, "."); i != -1 {
		return strings.ToLower(in[i+1:])
	}
	return in
}

func wrapCommentWithLength(in string, length int) string {
	return txt.Wrap("// ", in, length)
}

func add(left, right int) int {
	return left + right
}

func sub(left, right int) int {
	return left - right
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
			if strings.HasSuffix(name, "_gen.go") || strings.HasSuffix(name, "_gen_test.go") {
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
	// I intentionally do not include the standard "Code generated ... DO NOT EDIT" comment at the top of the file
	// because I want linting to run on the file. Instead, this comment is inserted, which is similar, but doesn't
	// trigger the linter to ignore the file.
	fmt.Fprintf(&buffer, "// Code created from %q - don't edit by hand\n//\n", tmpl)
	if err := tmpls.ExecuteTemplate(&buffer, tmpl, tmplArg); err != nil {
		return errs.Wrap(err)
	}
	data := buffer.Bytes()
	if formatted, err := format.Source(data); err != nil {
		fmt.Println(errs.NewWithCause(fmt.Sprintf("unable to format %q", outputPath), err))
	} else {
		data = formatted
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o750); err != nil {
		return errs.Wrap(err)
	}
	if err := os.WriteFile(outputPath, data, 0o640); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

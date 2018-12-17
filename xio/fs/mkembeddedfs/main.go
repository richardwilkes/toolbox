package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/collection"
	"github.com/richardwilkes/toolbox/txt"
)

type data struct {
	FSPath     string
	Path       string
	Name       string
	Size       int64
	ModTime    int64
	Compressed bool
	Data       string
}

type tmplInput struct {
	Tag            string
	Pkg            string
	Var            string
	Files          []*data
	NoEmbedModTime bool
}

func main() {
	cmdline.AppName = "Make Embedded Go FileSystem"
	cmdline.AppVersion = "1.0"
	cmdline.CopyrightYears = "2018"
	cmdline.CopyrightHolder = "Richard A. Wilkes"
	cl := cmdline.New(true)
	cl.UsageSuffix = "<one or more file paths to include>"
	cfg := tmplInput{
		Pkg: "main",
		Var: "EFS",
	}
	var strip, ignore string
	var output = "efs.go"
	cl.NewStringOption(&cfg.Pkg).SetSingle('p').SetName("pkg").SetUsage("The package name for the output file")
	cl.NewStringOption(&strip).SetSingle('s').SetName("strip").SetUsage("A prefix to remove from stored file paths")
	cl.NewStringOption(&ignore).SetSingle('i').SetName("ignore").SetUsage("A regular expression for file paths to ignore")
	cl.NewStringOption(&output).SetSingle('o').SetName("output").SetUsage("The output file path")
	cl.NewStringOption(&cfg.Var).SetSingle('n').SetName("name").SetUsage("The variable name to use for the embedded filesystem")
	cl.NewStringOption(&cfg.Tag).SetSingle('t').SetName("tag").SetUsage("A build tag to guard the output file with")
	cl.NewBoolOption(&cfg.NoEmbedModTime).SetSingle('N').SetName("no-modtime").SetUsage("Don't embed the modification time of files")
	paths := cl.Parse(os.Args[1:])
	if len(paths) == 0 {
		fail("Must specify at least one input path to process")
	}
	if output == "" {
		fail("The output file path may not be empty")
	}
	if cfg.Var == "" {
		fail("The variable name may not be empty")
	}
	c := collector{paths: collection.NewStringSet()}
	if ignore != "" {
		var err error
		c.ignoreRegex, err = regexp.Compile(ignore)
		failIfErr(err)
	}
	for _, one := range paths {
		failIfErr(filepath.Walk(filepath.Clean(one), c.walk))
	}
	all, err := c.prepare(strip)
	failIfErr(err)
	cfg.Files = all

	tmpl := template.Must(template.New("").Parse(pkgTemplate))

	var buffer bytes.Buffer
	failIfErr(tmpl.Execute(&buffer, &cfg))
	d, err := format.Source(buffer.Bytes())
	failIfErr(err)
	failIfErr(ioutil.WriteFile(output, d, 0644))
	atexit.Exit(0)
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	atexit.Exit(1)
}

func failIfErr(err error) {
	if err != nil {
		fail(err.Error())
	}
}

type collector struct {
	ignoreRegex *regexp.Regexp
	paths       collection.StringSet
}

func (c *collector) walk(path string, info os.FileInfo, err error) error {
	if info == nil {
		if err == nil {
			err = fmt.Errorf("Unable to walk to %q", path)
		}
		return err
	}
	if c.ignoreRegex != nil && c.ignoreRegex.MatchString(path) {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}
	name := info.Name()
	if name == ".git" || name == ".svn" || name == ".DS_Store" {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}
	if !info.IsDir() {
		c.paths.Add(filepath.ToSlash(path))
	}
	return nil
}

func (c *collector) prepare(strip string) ([]*data, error) {
	var all []*data
	paths := collection.NewStringSet()
	for _, one := range c.paths.Values() {
		fsPath := one
		one = strings.TrimPrefix(one, strip)
		if paths.Contains(one) {
			return nil, errors.New("When prefix is stripped, more than one file maps to the same path: " + one)
		}
		paths.Add(one)
		all = append(all, &data{
			FSPath: fsPath,
			Path:   path.Clean("/" + one),
		})
	}
	sort.Slice(all, func(i, j int) bool { return txt.NaturalLess(all[i].Path, all[j].Path, false) })
	for _, one := range all {
		f, err := os.Open(one.FSPath)
		failIfErr(err)
		fi, err := f.Stat()
		failIfErr(err)
		one.Name = fi.Name()
		one.Size = fi.Size()
		one.ModTime = fi.ModTime().UnixNano()
		var raw bytes.Buffer
		_, err = raw.ReadFrom(f)
		failIfErr(err)
		failIfErr(f.Close())
		var compressed bytes.Buffer
		gw := gzip.NewWriter(&compressed)
		_, err = io.Copy(gw, bytes.NewBuffer(raw.Bytes()))
		failIfErr(err)
		failIfErr(gw.Close())
		var data []byte
		if compressed.Len() < raw.Len() {
			one.Compressed = true
			data = compressed.Bytes()
		} else {
			data = raw.Bytes()
		}
		var buffer bytes.Buffer
		count := 0
		for i := 0; i < len(data); i++ {
			switch count {
			case 0:
			case 16:
				buffer.WriteString("\n\t\t")
				count = 0
			default:
				buffer.WriteByte(' ')
			}
			fmt.Fprintf(&buffer, "0x%02x,", data[i])
			count++
		}
		one.Data = buffer.String()
	}
	return all, nil
}

var pkgTemplate = `// Code generated - DO NOT EDIT.
{{if .Tag}}
// {{/**/}}+build {{.Tag}}
{{end}}{{$nomodtime := .NoEmbedModTime}}
package {{.Pkg}}

import (
	"time"

	"github.com/richardwilkes/toolbox/xio/fs/embedded"
)

// {{.Var}} holds an embedded filesystem.
var {{.Var}} = embedded.NewEFS(map[string]*embedded.File{
{{- range .Files}}
	{{printf "%q" .Path}}: embedded.NewFile({{printf "%q" .Name}}, {{if $nomodtime}}time.Now(){{else}}time.Unix(0, {{.ModTime}}){{end}}, {{.Size}}, {{.Compressed}}, []byte{
		{{.Data}}
	}),
{{- end}}
})
`

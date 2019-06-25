package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/i18n"
)

func main() {
	cmdline.CopyrightYears = "2016-2017"
	cmdline.CopyrightHolder = "Richard A. Wilkes"
	cmdline.License = "Mozilla Public License 2.0"
	cl := cmdline.New(true)
	cl.UsageSuffix = "<path> [path...]"
	cl.Description = i18n.Text("Generates a template for a localization file from source code.")
	outPath := "language.i18n"
	cl.NewStringOption(&outPath).SetSingle('o').SetName("output").SetArg("path").SetUsage("The output file")
	args := cl.Parse(os.Args[1:])
	if outPath == "" {
		cl.FatalMsg(i18n.Text("The output file may not be an empty path."))
	}
	if len(args) == 0 {
		cl.FatalMsg(i18n.Text("At least one path must be specified."))
	}
	kv := make(map[string]string)
	fset := token.NewFileSet()
	for _, pathArg := range args {
		var err error
		if pathArg, err = filepath.Abs(pathArg); err == nil {
			walkErr := filepath.Walk(pathArg, func(path string, fi os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !fi.IsDir() && filepath.Ext(path) == ".go" {
					fmt.Println(path)
					var file *ast.File
					if file, err = parser.ParseFile(fset, path, nil, 0); err != nil {
						fmt.Fprintln(os.Stderr, err)
						atexit.Exit(1)
					}
					const (
						LookForPackageState = iota
						LookForTextCallState
						LookForParameterState
					)
					state := LookForPackageState
					ast.Inspect(file, func(node ast.Node) bool {
						switch x := node.(type) {
						case *ast.Ident:
							switch state {
							case LookForPackageState:
								if x.Name == "i18n" {
									state = LookForTextCallState
								}
							case LookForTextCallState:
								if x.Name == "Text" {
									state = LookForParameterState
								} else {
									state = LookForPackageState
								}
							default:
								state = LookForPackageState
							}
						case *ast.BasicLit:
							if state == LookForParameterState {
								if x.Kind == token.STRING {
									kv[x.Value] = x.Value
								}
							}
							state = LookForPackageState
						case nil:
						default:
							state = LookForPackageState
						}
						return true
					})
				}
				return nil
			})
			if walkErr != nil {
				fmt.Fprintln(os.Stderr, walkErr)
			}
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	out, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create '%s'.\n", outPath)
		atexit.Exit(1)
	}
	fmt.Fprintf(out, "# Generated on %v\n", time.Now())
	keys := make([]string, 0, len(kv))
	for key := range kv {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(out, "\nk: %s\nv: %s\n", key, key)
	}
	if err := out.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	atexit.Exit(0)
}

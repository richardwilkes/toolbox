// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/v2/i18n"
	"github.com/richardwilkes/toolbox/v2/xflag"
	"github.com/richardwilkes/toolbox/v2/xos"
	"github.com/richardwilkes/toolbox/v2/xstrings"
)

func main() {
	xos.AppIdentifier = "com.trollworks.i18n"
	xos.CopyrightStartYear = "2016"
	xos.CopyrightHolder = "Richard A. Wilkes"
	xos.License = "Mozilla Public License 2.0"
	xflag.SetUsage(nil, i18n.Text("Generates a template for a localization file from source code."), "<path> [path...]")
	outPath := flag.String("output", "language.i18n", "The output `path`")
	xflag.Parse()
	if *outPath == "" {
		xos.ExitWithMsg(i18n.Text("The output file may not be an empty path."))
	}
	args := flag.Args()
	if len(args) == 0 {
		xos.ExitWithMsg(i18n.Text("At least one path must be specified."))
	}
	kv := make(map[string]string)
	fileSet := token.NewFileSet()
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
					if file, err = parser.ParseFile(fileSet, path, nil, 0); err != nil {
						fmt.Fprintln(os.Stderr, err)
						xos.Exit(1)
					}
					ast.Inspect(file, func(node ast.Node) bool {
						call, ok := node.(*ast.CallExpr)
						if !ok || len(call.Args) == 0 {
							return true
						}
						var sel *ast.SelectorExpr
						if sel, ok = call.Fun.(*ast.SelectorExpr); !ok || sel.Sel.Name != "Text" {
							return true
						}
						var pkg *ast.Ident
						if pkg, ok = sel.X.(*ast.Ident); !ok || pkg.Name != "i18n" {
							return true
						}
						// Only the first argument is the text to localize. Any arguments that follow are values to
						// be formatted into it, so they must not be extracted, even if they are string literals.
						var lit *ast.BasicLit
						if lit, ok = call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
							var v string
							if v, err = strconv.Unquote(lit.Value); err != nil {
								fmt.Fprintln(os.Stderr, err)
							} else {
								kv[v] = v
							}
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

	keys := make([]string, 0, len(kv))
	for key := range kv {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return xstrings.NaturalLess(keys[i], keys[j], true)
	})
	out, err := os.OpenFile(*outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create '%s'.\n", *outPath)
		xos.Exit(1)
	}
	fmt.Fprintf(out, `# Generated on %v
#
# Key-value pairs are defined as one or more lines prefixed with "k:" for the
# key, followed by one or more lines prefixed with "v:" for the value. These
# prefixes are then followed by a quoted string, using escaping rules for Go
# strings where needed. When two or more lines are present in a row, they will
# be concatenated together with an intervening \n character.
#
# Do NOT modify the 'k' values. They are the values as seen in the code.
#
# Replace the 'v' values with the appropriate translation.
`, time.Now().Format(time.RFC1123))
	for _, key := range keys {
		fmt.Fprintln(out)
		for p := range strings.SplitSeq(key, "\n") {
			if _, err = fmt.Fprintf(out, "k:%q\n", p); err != nil {
				fmt.Fprintln(os.Stderr, err)
				xos.Exit(1)
			}
		}
		for p := range strings.SplitSeq(key, "\n") {
			if _, err = fmt.Fprintf(out, "v:%q\n", p); err != nil {
				fmt.Fprintln(os.Stderr, err)
				xos.Exit(1)
			}
		}
	}
	if err = out.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		xos.Exit(1)
	}
	xos.Exit(0)
}

/*
 * Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

// Package xlsx provides the ability to extract text from Excel spreadsheets.
package xlsx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xio"
)

// Sheet holds the data contained in a single worksheet.
type Sheet struct {
	Name  string
	Min   Ref
	Max   Ref
	Cells map[Ref]Cell
}

// Load sheets from an .xlsx file.
func Load(path string) ([]Sheet, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(r)
	return load(&r.Reader)
}

// Read sheets from an .xlsx stream.
func Read(in io.ReaderAt, size int64) ([]Sheet, error) {
	r, err := zip.NewReader(in, size)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return load(r)
}

func load(r *zip.Reader) ([]Sheet, error) {
	var sheetNames []string
	var strs []string
	var files []*zip.File
	var err error
	for _, f := range r.File {
		switch {
		case f.Name == "docProps/app.xml":
			if sheetNames, err = loadSheetNames(f); err != nil {
				return nil, err
			}
		case f.Name == "xl/sharedStrings.xml":
			if strs, err = loadStrings(f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.Name, "xl/worksheets/sheet"):
			files = append(files, f)
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return txt.NaturalLess(files[i].Name, files[j].Name, true)
	})
	sheets := make([]Sheet, 0, len(files))
	for i, f := range files {
		var sheet *Sheet
		if sheet, err = loadSheet(f, strs); err != nil {
			return nil, err
		}
		if i < len(sheetNames) {
			sheet.Name = sheetNames[i]
		} else {
			sheet.Name = fmt.Sprintf("Sheet%d", i+1)
		}
		sheets = append(sheets, *sheet)
	}
	return sheets, nil
}

func loadSheetNames(f *zip.File) ([]string, error) {
	fr, err := f.Open()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(fr)
	decoder := xml.NewDecoder(fr)
	var data struct {
		Names []string `xml:"TitlesOfParts>vector>lpstr"`
	}
	if err = decoder.Decode(&data); err != nil {
		return nil, errs.Wrap(err)
	}
	return data.Names, nil
}

func loadStrings(f *zip.File) ([]string, error) {
	fr, err := f.Open()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(fr)
	decoder := xml.NewDecoder(fr)
	var data struct {
		SST []string `xml:"si>t"`
	}
	if err = decoder.Decode(&data); err != nil {
		return nil, errs.Wrap(err)
	}
	return data.SST, nil
}

func loadSheet(f *zip.File, strs []string) (*Sheet, error) {
	fr, err := f.Open()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(fr)
	decoder := xml.NewDecoder(fr)
	var data struct {
		Cells []struct {
			Label string  `xml:"r,attr"`
			Type  string  `xml:"t,attr"`
			Value *string `xml:"v"`
		} `xml:"sheetData>row>c"`
	}
	if err = decoder.Decode(&data); err != nil {
		return nil, errs.Wrap(err)
	}
	sheet := &Sheet{
		Min:   Ref{Row: math.MaxInt32, Col: math.MaxInt32},
		Max:   Ref{},
		Cells: make(map[Ref]Cell, len(data.Cells)),
	}
	for _, one := range data.Cells {
		if one.Value == nil {
			continue
		}
		ref := ParseRef(one.Label)
		cell := Cell{Value: *one.Value}
		switch one.Type {
		case "s": // String
			var v int
			if v, err = strconv.Atoi(cell.Value); err != nil {
				return nil, errs.Wrap(err)
			}
			if v >= 0 && v < len(strs) {
				cell.Value = strs[v]
			} else {
				cell.Value = "#REF!"
			}
			cell.Type = String
		case "b": // Boolean
			cell.Type = Boolean
		default: // Number
			cell.Type = Number
		}
		if sheet.Min.Row > ref.Row {
			sheet.Min.Row = ref.Row
		}
		if sheet.Min.Col > ref.Col {
			sheet.Min.Col = ref.Col
		}
		if sheet.Max.Row < ref.Row {
			sheet.Max.Row = ref.Row
		}
		if sheet.Max.Col < ref.Col {
			sheet.Max.Col = ref.Col
		}
		sheet.Cells[ref] = cell
	}
	if sheet.Min.Row > sheet.Max.Row {
		sheet.Min.Row = sheet.Max.Row
	}
	if sheet.Min.Col > sheet.Max.Col {
		sheet.Min.Col = sheet.Max.Col
	}
	return sheet, nil
}

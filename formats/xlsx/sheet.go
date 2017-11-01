package xlsx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/richardwilkes/gokit/errs"
	"github.com/richardwilkes/gokit/txt"
	"github.com/richardwilkes/gokit/xio"
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
	var sheets []Sheet
	var sheetNames []string
	var strs []string
	var files []*zip.File
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
	for i, f := range files {
		sheet, err := loadSheet(f, strs)
		if err != nil {
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
		Cells: make(map[Ref]Cell),
	}
	for _, one := range data.Cells {
		if one.Value != nil {
			ref := ParseRef(one.Label)
			cell := Cell{Value: *one.Value}
			switch one.Type {
			case "s": // String
				v, err := strconv.Atoi(cell.Value)
				if err != nil {
					return nil, errs.Wrap(err)
				}
				if v < 0 || v >= len(strs) {
					return nil, errs.New("String index out of bounds")
				}
				cell.Type = String
				cell.Value = strs[v]
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
	}
	if sheet.Min.Row > sheet.Max.Row {
		sheet.Min.Row = sheet.Max.Row
	}
	if sheet.Min.Col > sheet.Max.Col {
		sheet.Min.Col = sheet.Max.Col
	}
	return sheet, nil
}

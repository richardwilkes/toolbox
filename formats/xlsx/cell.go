package xlsx

import (
	"strconv"
	"time"
)

// Cell types.
const (
	String CellType = iota
	Number
	Boolean
)

// CellType holds an enumeration of cell types.
type CellType int

// Cell holds the contents of a cell.
type Cell struct {
	Type  CellType
	Value string
}

func (c *Cell) String() string {
	return c.Value
}

// Integer returns the value of this cell as an integer.
func (c *Cell) Integer() int {
	v, err := strconv.Atoi(c.Value)
	if err != nil {
		v = int(c.Float())
	}
	return v
}

// Float returns the value of this cell as an float.
func (c *Cell) Float() float64 {
	f, err := strconv.ParseFloat(c.Value, 64)
	if err != nil {
		return 0
	}
	return f
}

// Boolean returns the value of this cell as a boolean.
func (c *Cell) Boolean() bool {
	return c.Value != "0"
}

// Time returns the value of this cell as a time.Time.
func (c *Cell) Time() time.Time {
	return timeFromExcelTime(c.Float())
}

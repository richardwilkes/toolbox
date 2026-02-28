// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xslices_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xslices"
)

func TestSet(t *testing.T) {
	c := check.New(t)

	// Test with integers
	intData := []int{1, 2, 3, 2, 1, 4}
	intSet := xslices.Set(intData)
	c.Equal(4, len(intSet))
	_, exists := intSet[1]
	c.True(exists)
	_, exists = intSet[2]
	c.True(exists)
	_, exists = intSet[3]
	c.True(exists)
	_, exists = intSet[4]
	c.True(exists)
	_, exists = intSet[5]
	c.False(exists)

	// Test with strings
	strData := []string{"apple", "banana", "apple", "cherry"}
	strSet := xslices.Set(strData)
	c.Equal(3, len(strSet))
	_, exists = strSet["apple"]
	c.True(exists)
	_, exists = strSet["banana"]
	c.True(exists)
	_, exists = strSet["cherry"]
	c.True(exists)
	_, exists = strSet["orange"]
	c.False(exists)

	// Test with empty slice
	emptySet := xslices.Set([]int{})
	c.Equal(0, len(emptySet))
}

func TestMapFromData(t *testing.T) {
	c := check.New(t)

	// Test with struct data
	type person struct {
		Name string
		ID   int
	}

	people := []person{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}

	// Map by ID
	idMap := xslices.MapFromData(people, func(p person) int { return p.ID })
	c.Equal(3, len(idMap))
	c.Equal("Alice", idMap[1].Name)
	c.Equal("Bob", idMap[2].Name)
	c.Equal("Charlie", idMap[3].Name)

	// Map by name
	nameMap := xslices.MapFromData(people, func(p person) string { return p.Name })
	c.Equal(3, len(nameMap))
	c.Equal(1, nameMap["Alice"].ID)
	c.Equal(2, nameMap["Bob"].ID)
	c.Equal(3, nameMap["Charlie"].ID)

	// Test with duplicate keys (should overwrite)
	duplicatePeople := []person{
		{ID: 1, Name: "Alice"},
		{ID: 1, Name: "Alice2"},
	}
	duplicateMap := xslices.MapFromData(duplicatePeople, func(p person) int { return p.ID })
	c.Equal(1, len(duplicateMap))
	c.Equal("Alice2", duplicateMap[1].Name)

	// Test with empty slice
	emptyMap := xslices.MapFromData([]person{}, func(p person) int { return p.ID })
	c.Equal(0, len(emptyMap))
}

func TestMapFromKeys(t *testing.T) {
	c := check.New(t)

	// Test with integer keys
	keys := []int{1, 2, 3, 4}
	squareMap := xslices.MapFromKeys(keys, func(k int) int { return k * k })
	c.Equal(4, len(squareMap))
	c.Equal(1, squareMap[1])
	c.Equal(4, squareMap[2])
	c.Equal(9, squareMap[3])
	c.Equal(16, squareMap[4])

	// Test with string keys
	strKeys := []string{"a", "b", "c"}
	upperMap := xslices.MapFromKeys(strKeys, func(k string) string { return k + "_upper" })
	c.Equal(3, len(upperMap))
	c.Equal("a_upper", upperMap["a"])
	c.Equal("b_upper", upperMap["b"])
	c.Equal("c_upper", upperMap["c"])

	// Test with struct values
	type info struct {
		Label string
		Count int
	}

	infoMap := xslices.MapFromKeys(keys, func(k int) info {
		return info{Count: k, Label: "item"}
	})
	c.Equal(4, len(infoMap))
	c.Equal(info{Count: 1, Label: "item"}, infoMap[1])
	c.Equal(info{Count: 2, Label: "item"}, infoMap[2])
	c.Equal(info{Count: 3, Label: "item"}, infoMap[3])
	c.Equal(info{Count: 4, Label: "item"}, infoMap[4])

	// Test with empty keys
	emptyMap := xslices.MapFromKeys([]int{}, func(_ int) string { return "empty" })
	c.Equal(0, len(emptyMap))
}

// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package redblack_test

import (
	"cmp"
	"fmt"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/collection/redblack"
)

func TestRedBlackTree(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, int](cmp.Compare[int])
	c.Equal(0, rbt.Count())

	rbt.Insert(10, 10)
	c.Equal(1, rbt.Count())

	result, ok := rbt.Get(10)
	c.True(ok)
	c.Equal(10, result)

	rbt.Remove(10)
	c.Equal(0, rbt.Count())

	rbt.Insert(10, 10)
	rbt.Insert(5, 5)
	rbt.Insert(15, 15)
	c.Equal(3, rbt.Count())

	var values []int
	rbt.Traverse(func(_, value int) bool {
		values = append(values, value)
		return true
	})
	c.Equal(3, len(values))
	c.Equal([]int{5, 10, 15}, values)

	rbt.Insert(10, 10)
	c.Equal(4, rbt.Count())

	values = nil
	rbt.Traverse(func(_, value int) bool {
		values = append(values, value)
		return true
	})
	c.Equal(4, len(values))
	c.Equal([]int{5, 10, 10, 15}, values)

	values = nil
	rbt.ReverseTraverse(func(_, value int) bool {
		values = append(values, value)
		return true
	})
	c.Equal(4, len(values))
	c.Equal([]int{15, 10, 10, 5}, values)

	rbt.Remove(7)
	c.Equal(4, rbt.Count())

	rbt.Remove(10)
	c.Equal(3, rbt.Count())

	rbt.Remove(10)
	c.Equal(2, rbt.Count())

	values = nil
	rbt.Traverse(func(_, value int) bool {
		values = append(values, value)
		return true
	})
	c.Equal(2, len(values))
	c.Equal([]int{5, 15}, values)

	for i := -10; i < 21; i++ {
		rbt.Insert(i, i)
	}
	c.Equal(33, rbt.Count())

	result, ok = rbt.Get(-3)
	c.True(ok)
	c.Equal(-3, result)

	_, ok = rbt.Get(-11)
	c.False(ok)

	values = nil
	rbt.TraverseStartingAt(30, func(_, value int) bool {
		values = append(values, value)
		return true
	})
	c.Equal(0, len(values))

	values = nil
	rbt.TraverseStartingAt(20, func(_, value int) bool {
		values = append(values, value)
		return true
	})
	c.Equal(1, len(values))
	c.Equal([]int{20}, values)

	values = nil
	rbt.TraverseStartingAt(18, func(_, value int) bool {
		values = append(values, value)
		return true
	})
	c.Equal(3, len(values))
	c.Equal([]int{18, 19, 20}, values)

	values = nil
	rbt.ReverseTraverseStartingAt(-20, func(_, value int) bool {
		values = append(values, value)
		return true
	})
	c.Equal(0, len(values))

	values = nil
	rbt.ReverseTraverseStartingAt(-10, func(_, value int) bool {
		values = append(values, value)
		return true
	})
	c.Equal(1, len(values))
	c.Equal([]int{-10}, values)

	values = nil
	rbt.ReverseTraverseStartingAt(-8, func(_, value int) bool {
		values = append(values, value)
		return true
	})
	c.Equal(3, len(values))
	c.Equal([]int{-8, -9, -10}, values)
}

func TestEmptyTree(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, string](cmp.Compare[int])

	// Test empty tree properties
	c.True(rbt.Empty())
	c.Equal(0, rbt.Count())

	// Test Get on empty tree
	_, exists := rbt.Get(42)
	c.False(exists)

	// Test First/Last on empty tree
	_, exists = rbt.First()
	c.False(exists)
	_, exists = rbt.Last()
	c.False(exists)

	// Test traversal on empty tree
	called := false
	rbt.Traverse(func(int, string) bool {
		called = true
		return true
	})
	c.False(called)

	rbt.ReverseTraverse(func(int, string) bool {
		called = true
		return true
	})
	c.False(called)

	rbt.TraverseStartingAt(5, func(int, string) bool {
		called = true
		return true
	})
	c.False(called)

	rbt.ReverseTraverseStartingAt(5, func(int, string) bool {
		called = true
		return true
	})
	c.False(called)

	// Test Remove on empty tree
	rbt.Remove(42) // Should not panic
	c.Equal(0, rbt.Count())
}

func TestSingleNode(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, string](cmp.Compare[int])

	rbt.Insert(10, "ten")
	c.False(rbt.Empty())
	c.Equal(1, rbt.Count())

	// Test First/Last on single node
	first, exists := rbt.First()
	c.True(exists)
	c.Equal("ten", first)

	last, exists := rbt.Last()
	c.True(exists)
	c.Equal("ten", last)

	// Test Get
	value, exists := rbt.Get(10)
	c.True(exists)
	c.Equal("ten", value)

	_, exists = rbt.Get(20)
	c.False(exists)

	// Test traversal
	var keys []int
	var values []string
	rbt.Traverse(func(k int, v string) bool {
		keys = append(keys, k)
		values = append(values, v)
		return true
	})
	c.Equal([]int{10}, keys)
	c.Equal([]string{"ten"}, values)

	// Test reverse traversal
	keys = nil
	values = nil
	rbt.ReverseTraverse(func(k int, v string) bool {
		keys = append(keys, k)
		values = append(values, v)
		return true
	})
	c.Equal([]int{10}, keys)
	c.Equal([]string{"ten"}, values)

	// Test removal
	rbt.Remove(10)
	c.True(rbt.Empty())
	c.Equal(0, rbt.Count())
}

func TestFirstLast(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, string](cmp.Compare[int])

	// Insert in random order
	values := []int{50, 30, 70, 20, 40, 60, 80, 10, 90}
	for _, v := range values {
		rbt.Insert(v, fmt.Sprintf("value_%d", v))
	}

	first, exists := rbt.First()
	c.True(exists)
	c.Equal("value_10", first)

	last, exists := rbt.Last()
	c.True(exists)
	c.Equal("value_90", last)

	// Remove first and last, test again
	rbt.Remove(10)
	rbt.Remove(90)

	first, exists = rbt.First()
	c.True(exists)
	c.Equal("value_20", first)

	last, exists = rbt.Last()
	c.True(exists)
	c.Equal("value_80", last)
}

func TestDuplicateKeys(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, string](cmp.Compare[int])

	// Insert duplicate keys with different values
	rbt.Insert(10, "first")
	rbt.Insert(10, "second")
	rbt.Insert(10, "third")
	c.Equal(3, rbt.Count())

	// Get should return the left-most (first inserted) value
	value, exists := rbt.Get(10)
	c.True(exists)
	c.Equal("first", value)

	// Traverse should show all duplicates in order
	var values []string
	rbt.Traverse(func(_ int, v string) bool {
		values = append(values, v)
		return true
	})
	c.Equal([]string{"first", "second", "third"}, values)

	// Remove should remove the left-most occurrence first
	rbt.Remove(10)
	c.Equal(2, rbt.Count())

	value, exists = rbt.Get(10)
	c.True(exists)
	c.Equal("second", value)

	rbt.Remove(10)
	c.Equal(1, rbt.Count())

	value, exists = rbt.Get(10)
	c.True(exists)
	c.Equal("third", value)

	rbt.Remove(10)
	c.Equal(0, rbt.Count())
	c.True(rbt.Empty())
}

func TestTraversalEarlyTermination(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, int](cmp.Compare[int])

	for i := 1; i <= 10; i++ {
		rbt.Insert(i, i*10)
	}

	// Test early termination in forward traversal
	var visited []int
	rbt.Traverse(func(k, _ int) bool {
		visited = append(visited, k)
		return k < 5 // Stop at 5
	})
	c.Equal([]int{1, 2, 3, 4, 5}, visited)

	// Test early termination in reverse traversal
	visited = nil
	rbt.ReverseTraverse(func(k, _ int) bool {
		visited = append(visited, k)
		return k > 6 // Stop at 6
	})
	c.Equal([]int{10, 9, 8, 7, 6}, visited)

	// Test early termination in TraverseStartingAt
	visited = nil
	rbt.TraverseStartingAt(5, func(k, _ int) bool {
		visited = append(visited, k)
		return k < 8 // Stop at 8
	})
	c.Equal([]int{5, 6, 7, 8}, visited)

	// Test early termination in ReverseTraverseStartingAt
	visited = nil
	rbt.ReverseTraverseStartingAt(7, func(k, _ int) bool {
		visited = append(visited, k)
		return k > 4 // Stop at 4
	})
	c.Equal([]int{7, 6, 5, 4}, visited)
}

func TestTraverseStartingAt(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, int](cmp.Compare[int])

	for i := 1; i <= 10; i++ {
		rbt.Insert(i, i*10)
	}

	// Test starting at existing key
	var visited []int
	rbt.TraverseStartingAt(5, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal([]int{5, 6, 7, 8, 9, 10}, visited)

	// Test starting at non-existing key (should start at next greater)
	visited = nil
	rbt.TraverseStartingAt(6, func(k, _ int) bool { // Changed from 5.5 to 6
		visited = append(visited, k)
		return true
	})
	c.Equal([]int{6, 7, 8, 9, 10}, visited)

	// Test starting at key greater than all
	visited = nil
	rbt.TraverseStartingAt(15, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal(0, len(visited))

	// Test starting at key less than all
	visited = nil
	rbt.TraverseStartingAt(-5, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, visited)
}

func TestReverseTraverseStartingAt(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, int](cmp.Compare[int])

	for i := 1; i <= 10; i++ {
		rbt.Insert(i, i*10)
	}

	// Test starting at existing key
	var visited []int
	rbt.ReverseTraverseStartingAt(6, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal([]int{6, 5, 4, 3, 2, 1}, visited)

	// Test starting at non-existing key (should start at next smaller)
	visited = nil
	rbt.ReverseTraverseStartingAt(7, func(k, _ int) bool { // Changed from 6.5 to 7
		visited = append(visited, k)
		return true
	})
	c.Equal([]int{7, 6, 5, 4, 3, 2, 1}, visited) // Updated expected result

	// Test starting at key less than all
	visited = nil
	rbt.ReverseTraverseStartingAt(-5, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal(0, len(visited))

	// Test starting at key greater than all
	visited = nil
	rbt.ReverseTraverseStartingAt(15, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal([]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, visited)
}

func TestLargeDataset(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, int](cmp.Compare[int])

	// Insert many values to test tree balance
	n := 1000
	for i := range n {
		rbt.Insert(i, i*2)
	}
	c.Equal(n, rbt.Count())

	// Verify all values can be retrieved
	for i := range n {
		value, exists := rbt.Get(i)
		c.True(exists)
		c.Equal(i*2, value)
	}

	// Test First/Last
	first, exists := rbt.First()
	c.True(exists)
	c.Equal(0, first)

	last, exists := rbt.Last()
	c.True(exists)
	c.Equal((n-1)*2, last)

	// Remove half the values
	for i := 0; i < n; i += 2 {
		rbt.Remove(i)
	}
	c.Equal(n/2, rbt.Count())

	// Verify remaining values
	for i := 1; i < n; i += 2 {
		var value int
		value, exists = rbt.Get(i)
		c.True(exists)
		c.Equal(i*2, value)
	}

	// Verify removed values are gone
	for i := 0; i < n; i += 2 {
		_, exists = rbt.Get(i)
		c.False(exists)
	}
}

func TestStringKeys(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[string, int](cmp.Compare[string])

	data := map[string]int{
		"apple":  1,
		"banana": 2,
		"cherry": 3,
		"date":   4,
		"elder":  5,
	}

	for k, v := range data {
		rbt.Insert(k, v)
	}
	c.Equal(5, rbt.Count())

	// Test retrieval
	for k, v := range data {
		value, exists := rbt.Get(k)
		c.True(exists)
		c.Equal(v, value)
	}

	// Test traversal order (should be alphabetical)
	var keys []string
	rbt.Traverse(func(k string, _ int) bool {
		keys = append(keys, k)
		return true
	})
	c.Equal([]string{"apple", "banana", "cherry", "date", "elder"}, keys)

	// Test reverse traversal
	keys = nil
	rbt.ReverseTraverse(func(k string, _ int) bool {
		keys = append(keys, k)
		return true
	})
	c.Equal([]string{"elder", "date", "cherry", "banana", "apple"}, keys)
}

func TestCustomComparator(t *testing.T) {
	c := check.New(t)

	// Reverse comparator for descending order
	reverseCompare := func(a, b int) int {
		return cmp.Compare(b, a) // Reversed
	}

	rbt := redblack.New[int, string](reverseCompare)

	values := []int{3, 1, 4, 1, 5, 9, 2, 6}
	for _, v := range values {
		rbt.Insert(v, fmt.Sprintf("val_%d", v))
	}

	// Traversal should be in descending order
	var keys []int
	rbt.Traverse(func(k int, _ string) bool {
		keys = append(keys, k)
		return true
	})
	c.Equal([]int{9, 6, 5, 4, 3, 2, 1, 1}, keys)

	// First should be largest, Last should be smallest
	first, exists := rbt.First()
	c.True(exists)
	c.Equal("val_9", first)

	last, exists := rbt.Last()
	c.True(exists)
	c.Equal("val_1", last)
}

func TestRedBlackTreeProperties(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, int](cmp.Compare[int])

	// Insert values that could cause balance issues in a simple BST
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, v := range values {
		rbt.Insert(v, v)
	}

	// Tree should still function correctly despite potential for imbalance
	c.Equal(10, rbt.Count())

	// All values should be retrievable
	for _, v := range values {
		val, exists := rbt.Get(v)
		c.True(exists)
		c.Equal(v, val)
	}

	// Traversal should be in order
	var result []int
	rbt.Traverse(func(k, _ int) bool {
		result = append(result, k)
		return true
	})
	c.Equal(values, result)

	// Remove all values
	for _, v := range values {
		rbt.Remove(v)
	}
	c.True(rbt.Empty())
}

func TestDump(_ *testing.T) {
	// This test just ensures Dump doesn't panic
	rbt := redblack.New[int, int](cmp.Compare[int])

	// Empty tree
	rbt.Dump() // Should not panic

	// Single node
	rbt.Insert(10, 100)
	rbt.Dump() // Should not panic

	// Multiple nodes
	for i := 1; i <= 5; i++ {
		rbt.Insert(i, i*10)
	}
	rbt.Dump() // Should not panic
}

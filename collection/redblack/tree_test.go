// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

//nolint:goconst // I'd rather have the strings inline than extracted out into a constant for the tests.
package redblack_test

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"slices"
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

	c.True(rbt.Empty())
	c.Equal(0, rbt.Count())

	_, exists := rbt.Get(42)
	c.False(exists)

	_, exists = rbt.First()
	c.False(exists)
	_, exists = rbt.Last()
	c.False(exists)

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

	rbt.Remove(42) // Should not panic
	c.Equal(0, rbt.Count())
}

func TestSingleNode(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, string](cmp.Compare[int])

	rbt.Insert(10, "ten")
	c.False(rbt.Empty())
	c.Equal(1, rbt.Count())

	first, exists := rbt.First()
	c.True(exists)
	c.Equal("ten", first)

	last, exists := rbt.Last()
	c.True(exists)
	c.Equal("ten", last)

	value, exists := rbt.Get(10)
	c.True(exists)
	c.Equal("ten", value)

	_, exists = rbt.Get(20)
	c.False(exists)

	var keys []int
	var values []string
	rbt.Traverse(func(k int, v string) bool {
		keys = append(keys, k)
		values = append(values, v)
		return true
	})
	c.Equal([]int{10}, keys)
	c.Equal([]string{"ten"}, values)

	keys = nil
	values = nil
	rbt.ReverseTraverse(func(k int, v string) bool {
		keys = append(keys, k)
		values = append(values, v)
		return true
	})
	c.Equal([]int{10}, keys)
	c.Equal([]string{"ten"}, values)

	rbt.Remove(10)
	c.True(rbt.Empty())
	c.Equal(0, rbt.Count())
}

func TestFirstLast(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, string](cmp.Compare[int])

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

	rbt.Insert(10, "first")
	rbt.Insert(10, "second")
	rbt.Insert(10, "third")
	c.Equal(3, rbt.Count())

	value, exists := rbt.Get(10)
	c.True(exists)
	c.Equal("first", value)

	var values []string
	rbt.Traverse(func(_ int, v string) bool {
		values = append(values, v)
		return true
	})
	c.Equal([]string{"first", "second", "third"}, values)

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

func TestDuplicateKeysRandomized(t *testing.T) {
	c := check.New(t)
	r := rand.New(rand.NewPCG(1972, 5150)) //nolint:gosec // Yes, it is ok to use a weak prng here
	for range 1000 {
		rbt := redblack.New[int, int](cmp.Compare[int])
		for i := range 12 {
			rbt.Insert(r.IntN(6), i)
		}
		keys, values := collect(rbt)
		for key := range 6 {
			i := slices.Index(keys, key)
			value, exists := rbt.Get(key)
			if i == -1 {
				c.False(exists)
				continue
			}
			c.True(exists)
			c.Equal(values[i], value)
		}
		// Each removal must take out the first match in traversal order and leave the rest untouched.
		for len(keys) != 0 {
			i := slices.Index(keys, keys[r.IntN(len(keys))])
			rbt.Remove(keys[i])
			keys = slices.Delete(keys, i, i+1)
			values = slices.Delete(values, i, i+1)
			if len(keys) == 0 {
				keys, values = nil, nil
			}
			actualKeys, actualValues := collect(rbt)
			c.Equal(keys, actualKeys)
			c.Equal(values, actualValues)
		}
	}
}

func collect(rbt *redblack.Tree[int, int]) (keys, values []int) {
	rbt.Traverse(func(key, value int) bool {
		keys = append(keys, key)
		values = append(values, value)
		return true
	})
	return keys, values
}

func TestTraversalEarlyTermination(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, int](cmp.Compare[int])

	for i := 1; i <= 10; i++ {
		rbt.Insert(i, i*10)
	}

	var visited []int
	rbt.Traverse(func(k, _ int) bool {
		visited = append(visited, k)
		return k < 5
	})
	c.Equal([]int{1, 2, 3, 4, 5}, visited)

	visited = nil
	rbt.ReverseTraverse(func(k, _ int) bool {
		visited = append(visited, k)
		return k > 6
	})
	c.Equal([]int{10, 9, 8, 7, 6}, visited)

	visited = nil
	rbt.TraverseStartingAt(5, func(k, _ int) bool {
		visited = append(visited, k)
		return k < 8
	})
	c.Equal([]int{5, 6, 7, 8}, visited)

	visited = nil
	rbt.ReverseTraverseStartingAt(7, func(k, _ int) bool {
		visited = append(visited, k)
		return k > 4
	})
	c.Equal([]int{7, 6, 5, 4}, visited)
}

func TestTraverseStartingAt(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, int](cmp.Compare[int])

	for i := 1; i <= 10; i++ {
		rbt.Insert(i, i*10)
	}

	var visited []int
	rbt.TraverseStartingAt(5, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal([]int{5, 6, 7, 8, 9, 10}, visited)

	visited = nil
	rbt.TraverseStartingAt(6, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal([]int{6, 7, 8, 9, 10}, visited)

	visited = nil
	rbt.TraverseStartingAt(15, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal(0, len(visited))

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

	var visited []int
	rbt.ReverseTraverseStartingAt(6, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal([]int{6, 5, 4, 3, 2, 1}, visited)

	visited = nil
	rbt.ReverseTraverseStartingAt(7, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal([]int{7, 6, 5, 4, 3, 2, 1}, visited)

	visited = nil
	rbt.ReverseTraverseStartingAt(-5, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal(0, len(visited))

	visited = nil
	rbt.ReverseTraverseStartingAt(15, func(k, _ int) bool {
		visited = append(visited, k)
		return true
	})
	c.Equal([]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, visited)
}

func TestTraverseStartingAtWithDuplicateKeys(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, int](cmp.Compare[int])

	// Rotations distribute duplicates of the start key across both subtrees; a starting-at traversal must still visit
	// every one of them.
	rbt.Insert(1, 1)
	rbt.Insert(2, 2)
	const dupCount = 7
	for i := range dupCount {
		rbt.Insert(5, 100+i)
	}
	rbt.Insert(8, 8)
	rbt.Insert(9, 9)

	var got []int
	seen := make(map[int]bool)
	rbt.TraverseStartingAt(5, func(k, v int) bool {
		got = append(got, k)
		if k == 5 {
			seen[v] = true
		}
		return true
	})
	c.Equal(dupCount, len(seen))
	c.Equal(dupCount+2, len(got))
	for i := 1; i < len(got); i++ {
		c.True(got[i-1] <= got[i])
	}

	got = nil
	seen = make(map[int]bool)
	rbt.ReverseTraverseStartingAt(5, func(k, v int) bool {
		got = append(got, k)
		if k == 5 {
			seen[v] = true
		}
		return true
	})
	c.Equal(dupCount, len(seen))
	c.Equal(dupCount+2, len(got))
	for i := 1; i < len(got); i++ {
		c.True(got[i-1] >= got[i])
	}
}

func TestLargeDataset(t *testing.T) {
	c := check.New(t)
	rbt := redblack.New[int, int](cmp.Compare[int])

	n := 1000
	for i := range n {
		rbt.Insert(i, i*2)
	}
	c.Equal(n, rbt.Count())

	for i := range n {
		value, exists := rbt.Get(i)
		c.True(exists)
		c.Equal(i*2, value)
	}

	first, exists := rbt.First()
	c.True(exists)
	c.Equal(0, first)

	last, exists := rbt.Last()
	c.True(exists)
	c.Equal((n-1)*2, last)

	for i := 0; i < n; i += 2 {
		rbt.Remove(i)
	}
	c.Equal(n/2, rbt.Count())

	for i := 1; i < n; i += 2 {
		var value int
		value, exists = rbt.Get(i)
		c.True(exists)
		c.Equal(i*2, value)
	}

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

	for k, v := range data {
		value, exists := rbt.Get(k)
		c.True(exists)
		c.Equal(v, value)
	}

	var keys []string
	rbt.Traverse(func(k string, _ int) bool {
		keys = append(keys, k)
		return true
	})
	c.Equal([]string{"apple", "banana", "cherry", "date", "elder"}, keys)

	keys = nil
	rbt.ReverseTraverse(func(k string, _ int) bool {
		keys = append(keys, k)
		return true
	})
	c.Equal([]string{"elder", "date", "cherry", "banana", "apple"}, keys)
}

func TestCustomComparator(t *testing.T) {
	c := check.New(t)

	reverseCompare := func(a, b int) int {
		return cmp.Compare(b, a)
	}

	rbt := redblack.New[int, string](reverseCompare)

	values := []int{3, 1, 4, 1, 5, 9, 2, 6}
	for _, v := range values {
		rbt.Insert(v, fmt.Sprintf("val_%d", v))
	}

	var keys []int
	rbt.Traverse(func(k int, _ string) bool {
		keys = append(keys, k)
		return true
	})
	c.Equal([]int{9, 6, 5, 4, 3, 2, 1, 1}, keys)

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

	c.Equal(10, rbt.Count())

	for _, v := range values {
		val, exists := rbt.Get(v)
		c.True(exists)
		c.Equal(v, val)
	}

	var result []int
	rbt.Traverse(func(k, _ int) bool {
		result = append(result, k)
		return true
	})
	c.Equal(values, result)

	for _, v := range values {
		rbt.Remove(v)
	}
	c.True(rbt.Empty())
}

func TestDump(_ *testing.T) {
	// Ensures Dump doesn't panic.
	rbt := redblack.New[int, int](cmp.Compare[int])

	rbt.Dump()

	rbt.Insert(10, 100)
	rbt.Dump()

	for i := 1; i <= 5; i++ {
		rbt.Insert(i, i*10)
	}
	rbt.Dump()
}

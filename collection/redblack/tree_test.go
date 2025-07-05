// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
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

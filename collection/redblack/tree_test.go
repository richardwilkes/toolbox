// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/collection/redblack"
)

func TestRedBlackTree(t *testing.T) {
	rbt := redblack.New[int, int](cmp.Compare[int])
	check.Equal(t, 0, rbt.Count())

	rbt.Insert(10, 10)
	check.Equal(t, 1, rbt.Count())

	result, ok := rbt.Get(10)
	check.True(t, ok)
	check.Equal(t, 10, result)

	rbt.Remove(10)
	check.Equal(t, 0, rbt.Count())

	rbt.Insert(10, 10)
	rbt.Insert(5, 5)
	rbt.Insert(15, 15)
	check.Equal(t, 3, rbt.Count())

	var values []int
	rbt.Traverse(func(_, value int) bool {
		values = append(values, value)
		return true
	})
	check.Equal(t, 3, len(values))
	check.Equal(t, []int{5, 10, 15}, values)

	rbt.Insert(10, 10)
	check.Equal(t, 4, rbt.Count())

	values = nil
	rbt.Traverse(func(_, value int) bool {
		values = append(values, value)
		return true
	})
	check.Equal(t, 4, len(values))
	check.Equal(t, []int{5, 10, 10, 15}, values)

	values = nil
	rbt.ReverseTraverse(func(_, value int) bool {
		values = append(values, value)
		return true
	})
	check.Equal(t, 4, len(values))
	check.Equal(t, []int{15, 10, 10, 5}, values)

	rbt.Remove(7)
	check.Equal(t, 4, rbt.Count())

	rbt.Remove(10)
	check.Equal(t, 3, rbt.Count())

	rbt.Remove(10)
	check.Equal(t, 2, rbt.Count())

	values = nil
	rbt.Traverse(func(_, value int) bool {
		values = append(values, value)
		return true
	})
	check.Equal(t, 2, len(values))
	check.Equal(t, []int{5, 15}, values)

	for i := -10; i < 21; i++ {
		rbt.Insert(i, i)
	}
	check.Equal(t, 33, rbt.Count())

	result, ok = rbt.Get(-3)
	check.True(t, ok)
	check.Equal(t, -3, result)

	_, ok = rbt.Get(-11)
	check.False(t, ok)

	values = nil
	rbt.TraverseStartingAt(30, func(_, value int) bool {
		values = append(values, value)
		return true
	})
	check.Equal(t, 0, len(values))

	values = nil
	rbt.TraverseStartingAt(20, func(_, value int) bool {
		values = append(values, value)
		return true
	})
	check.Equal(t, 1, len(values))
	check.Equal(t, []int{20}, values)

	values = nil
	rbt.TraverseStartingAt(18, func(_, value int) bool {
		values = append(values, value)
		return true
	})
	check.Equal(t, 3, len(values))
	check.Equal(t, []int{18, 19, 20}, values)

	values = nil
	rbt.ReverseTraverseStartingAt(-20, func(_, value int) bool {
		values = append(values, value)
		return true
	})
	check.Equal(t, 0, len(values))

	values = nil
	rbt.ReverseTraverseStartingAt(-10, func(_, value int) bool {
		values = append(values, value)
		return true
	})
	check.Equal(t, 1, len(values))
	check.Equal(t, []int{-10}, values)

	values = nil
	rbt.ReverseTraverseStartingAt(-8, func(_, value int) bool {
		values = append(values, value)
		return true
	})
	check.Equal(t, 3, len(values))
	check.Equal(t, []int{-8, -9, -10}, values)
}

// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package redblack

import (
	"fmt"
	"strings"
)

type node[K, V any] struct {
	key    K
	value  V
	parent *node[K, V]
	left   *node[K, V]
	right  *node[K, V]
	black  bool
}

func (n *node[K, V]) isBlack() bool {
	return n == nil || n.black
}

func (n *node[K, V]) isRed() bool {
	return n != nil && !n.black
}

func (n *node[K, V]) find(compareFunc func(a, b K) int, key K) *node[K, V] {
	if n == nil {
		return nil
	}
	result := compareFunc(key, n.key)
	switch {
	case result < 0:
		return n.left.find(compareFunc, key)
	case result > 0:
		return n.right.find(compareFunc, key)
	default:
		// Always return the left-most one in the case of multiple matches
		cur := n
		for cur.left != nil && compareFunc(key, cur.left.key) == 0 {
			cur = cur.left
		}
		return cur
	}
}

func (n *node[K, V]) traverse(visitorFunc func(key K, value V) bool) bool {
	if n == nil {
		return true
	}
	if n.left.traverse(visitorFunc) {
		if visitorFunc(n.key, n.value) {
			return n.right.traverse(visitorFunc)
		}
	}
	return false
}

func (n *node[K, V]) reverseTraverse(visitorFunc func(key K, value V) bool) bool {
	if n == nil {
		return true
	}
	if n.right.reverseTraverse(visitorFunc) {
		if visitorFunc(n.key, n.value) {
			return n.left.reverseTraverse(visitorFunc)
		}
	}
	return false
}

func (n *node[K, V]) traverseEqualOrGreater(compareFunc func(a, b K) int, key K, visitorFunc func(key K, value V) bool) bool {
	if n == nil {
		return true
	}
	result := compareFunc(key, n.key)
	if result < 0 {
		if !n.left.traverseEqualOrGreater(compareFunc, key, visitorFunc) {
			return false
		}
	}
	if result <= 0 {
		if !visitorFunc(n.key, n.value) {
			return false
		}
	}
	return n.right.traverseEqualOrGreater(compareFunc, key, visitorFunc)
}

func (n *node[K, V]) traverseEqualOrLess(compareFunc func(a, b K) int, key K, visitorFunc func(key K, value V) bool) bool {
	if n == nil {
		return true
	}
	result := compareFunc(key, n.key)
	if result > 0 {
		if !n.right.traverseEqualOrLess(compareFunc, key, visitorFunc) {
			return false
		}
	}
	if result >= 0 {
		if !visitorFunc(n.key, n.value) {
			return false
		}
	}
	return n.left.traverseEqualOrLess(compareFunc, key, visitorFunc)
}

func (n *node[K, V]) dump(depth int, side string) {
	if n == nil {
		return
	}
	br := "r"
	if n.black {
		br = "b"
	}
	fmt.Printf("%s%s%s%v\n", strings.Repeat("  ", depth), br, side, n.key)
	n.left.dump(depth+1, "L ")
	n.right.dump(depth+1, "R ")
}

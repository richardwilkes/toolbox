// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package redblack

// Tree implements a Red-black Tree, as described here: https://en.wikipedia.org/wiki/Red–black_tree
type Tree[K, V any] struct {
	compare func(a, b K) int
	count   int
	root    *node[K, V]
}

// New creates a new red-black tree.
func New[K, V any](compareFunc func(a, b K) int) *Tree[K, V] {
	return &Tree[K, V]{compare: compareFunc}
}

// Empty returns true if the tree is empty.
func (t *Tree[K, V]) Empty() bool {
	return t.count == 0
}

// Count returns the number of nodes in the tree.
func (t *Tree[K, V]) Count() int {
	return t.count
}

// Get returns the first value that matches the given key.
func (t *Tree[K, V]) Get(key K) (value V, exists bool) {
	if n := t.root.find(t.compare, key); n != nil {
		return n.value, true
	}
	return value, false
}

// First returns the first value in the tree.
func (t *Tree[K, V]) First() (value V, exists bool) {
	n := t.root
	if n == nil {
		return value, false
	}
	for n.left != nil {
		n = n.left
	}
	return n.value, true
}

// Last returns the last value in the tree.
func (t *Tree[K, V]) Last() (value V, exists bool) {
	n := t.root
	if n == nil {
		return value, false
	}
	for n.right != nil {
		n = n.right
	}
	return n.value, true
}

// Dump a text version of the tree for debugging purposes.
func (t *Tree[K, V]) Dump() {
	t.root.dump(0, "")
}

// Traverse the tree, calling visitorFunc for each node, in order. If the visitorFunc returns false, the traversal will
// be aborted.
func (t *Tree[K, V]) Traverse(visitorFunc func(key K, value V) bool) {
	t.root.traverse(visitorFunc)
}

// TraverseStartingAt traverses the tree starting with the first node whose key is equal to or greater than the given
// key, calling visitorFunc for each node, in order. If the visitorFunc returns false, the traversal will be aborted.
func (t *Tree[K, V]) TraverseStartingAt(key K, visitorFunc func(key K, value V) bool) {
	t.root.traverseEqualOrGreater(t.compare, key, visitorFunc)
}

// ReverseTraverse traverses the tree, calling visitorFunc for each node, in reverse order. If the visitorFunc returns
// false, the traversal will be aborted.
func (t *Tree[K, V]) ReverseTraverse(visitorFunc func(key K, value V) bool) {
	t.root.reverseTraverse(visitorFunc)
}

// ReverseTraverseStartingAt traverses the tree starting with the last node whose key is equal to or less than the given
// key, calling visitorFunc for each node, in order. If the visitorFunc returns false, the traversal will be aborted.
func (t *Tree[K, V]) ReverseTraverseStartingAt(key K, visitorFunc func(key K, value V) bool) {
	t.root.traverseEqualOrLess(t.compare, key, visitorFunc)
}

// Insert a node into the tree.
func (t *Tree[K, V]) Insert(key K, value V) {
	n := &node[K, V]{key: key, value: value}
	cur := t.root
	n.parent = t.root
	for cur != nil {
		n.parent = cur
		if t.compare(key, cur.key) < 0 {
			cur = cur.left
		} else {
			cur = cur.right
		}
	}
	if n.parent == nil {
		t.root = n
	} else {
		if t.compare(key, n.parent.key) < 0 {
			n.parent.left = n
		} else {
			n.parent.right = n
		}
	}
	if n.parent != nil {
		parent := n.parent
		grandParent := parent.parent
		for grandParent != nil && parent.isRed() {
			if parent == grandParent.left {
				uncle := grandParent.right
				switch {
				case uncle.isRed():
					parent.black = true
					uncle.black = true
					grandParent.black = false
					n = grandParent
					parent = n.parent
					if parent != nil {
						grandParent = parent.parent
					} else {
						grandParent = nil
					}
				case n == parent.right:
					n, parent = parent, n
					t.rotateLeft(n)
				default:
					parent.black = true
					grandParent.black = false
					t.rotateRight(grandParent)
				}
			} else {
				uncle := grandParent.left
				switch {
				case uncle.isRed():
					parent.black = true
					uncle.black = true
					grandParent.black = false
					n = grandParent
					parent = n.parent
					if parent != nil {
						grandParent = parent.parent
					} else {
						grandParent = nil
					}
				case n == parent.left:
					n, parent = parent, n
					t.rotateRight(n)
				default:
					parent.black = true
					grandParent.black = false
					t.rotateLeft(grandParent)
				}
			}
		}
	}
	t.root.black = true
	t.count++
}

func (t *Tree[K, V]) rotateLeft(n *node[K, V]) {
	right := n.right
	n.right = right.left
	if right.left != nil {
		n.right.parent = n
	}
	right.parent = n.parent
	if n.parent != nil {
		if n.parent.left == n {
			n.parent.left = right
		} else {
			n.parent.right = right
		}
	} else {
		t.root = right
	}
	right.left = n
	n.parent = right
}

func (t *Tree[K, V]) rotateRight(n *node[K, V]) {
	left := n.left
	n.left = left.right
	if left.right != nil {
		n.left.parent = n
	}
	left.parent = n.parent
	if n.parent != nil {
		if n.parent.right == n {
			n.parent.right = left
		} else {
			n.parent.left = left
		}
	} else {
		t.root = left
	}
	left.right = n
	n.parent = left
}

// Remove a node from the tree. Note that if the key is not unique within the tree, the first key that matches on
// traversal will be chosen as the one to remove.
func (t *Tree[K, V]) Remove(key K) {
	n := t.root.find(t.compare, key)
	if n == nil {
		return
	}
	splice := n
	if n.left != nil && n.right != nil {
		splice = n.right
		for splice.left != nil {
			splice = splice.left
		}
	}
	var child *node[K, V]
	if splice.left != nil {
		child = splice.left
	} else {
		child = splice.right
	}
	if child != nil {
		child.parent = splice.parent
	}
	if splice.parent != nil {
		left := false
		parent := splice.parent
		if splice == parent.left {
			parent.left = child
			left = true
		} else {
			parent.right = child
		}
		if splice != n {
			n.key, splice.key = splice.key, n.key
			n.value, splice.value = splice.value, n.value
		}
		if splice.black {
			if child != nil {
				t.recolor(child)
			} else {
				child = splice
				child.parent = parent
				child.left = nil
				child.right = nil
				if left {
					parent.left = child
				} else {
					parent.right = child
				}
				t.recolor(child)
				if left {
					parent.left = nil
				} else {
					parent.right = nil
				}
			}
		}
	} else {
		t.root = child
	}
	if t.root != nil {
		t.root.black = true
	}
	t.count--
}

func (t *Tree[K, V]) recolor(n *node[K, V]) {
	for n != t.root && n.isBlack() {
		parent := n.parent
		switch {
		case parent.left == n:
			if sibling := parent.right; sibling != nil {
				if sibling.isRed() {
					sibling.black = true
					parent.black = false
					t.rotateLeft(parent)
					parent = n.parent
					sibling = parent.right
				}
				if sibling.left.isBlack() && sibling.right.isBlack() {
					sibling.black = false
					n = n.parent
				} else {
					if sibling.right.isBlack() {
						sibling.left.black = true
						sibling.black = false
						t.rotateRight(sibling)
						sibling = parent.right
					}
					sibling.black = parent.black
					parent.black = true
					sibling.right.black = true
					t.rotateLeft(parent)
					n = t.root
				}
			}
		case parent.right == n:
			if sibling := parent.left; sibling != nil {
				if sibling.isRed() {
					sibling.black = true
					parent.black = false
					t.rotateRight(parent)
					parent = n.parent
					sibling = parent.left
				}
				if sibling.right.isBlack() && sibling.left.isBlack() {
					sibling.black = false
					n = n.parent
				} else {
					if sibling.left.isBlack() {
						sibling.right.black = true
						sibling.black = false
						t.rotateLeft(sibling)
						sibling = parent.left
					}
					sibling.black = parent.black
					parent.black = true
					sibling.left.black = true
					t.rotateRight(parent)
					n = t.root
				}
			}
		default:
			parent.black = true
		}
	}
	n.black = true
}

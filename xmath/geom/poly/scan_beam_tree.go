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

package poly

import "golang.org/x/exp/constraints"

type scanBeamTree[T constraints.Float] struct {
	root    *scanBeamNode[T]
	entries int
}

type scanBeamNode[T constraints.Float] struct {
	y    T
	less *scanBeamNode[T]
	more *scanBeamNode[T]
}

func (s *scanBeamTree[T]) add(y T) {
	s.addToScanBeamTreeAt(&s.root, y)
}

func (s *scanBeamTree[T]) addToScanBeamTreeAt(node **scanBeamNode[T], y T) {
	switch {
	case *node == nil:
		*node = &scanBeamNode[T]{y: y}
		s.entries++
	case (*node).y > y:
		s.addToScanBeamTreeAt(&(*node).less, y)
	case (*node).y < y:
		s.addToScanBeamTreeAt(&(*node).more, y)
	default:
	}
}

func (s *scanBeamTree[T]) buildScanBeamTable() []T {
	table := make([]T, s.entries)
	if s.root != nil {
		s.root.buildScanBeamTableEntries(0, table)
	}
	return table
}

func (sbt *scanBeamNode[T]) buildScanBeamTableEntries(index int, table []T) int {
	if sbt.less != nil {
		index = sbt.less.buildScanBeamTableEntries(index, table)
	}
	table[index] = sbt.y
	index++
	if sbt.more != nil {
		index = sbt.more.buildScanBeamTableEntries(index, table)
	}
	return index
}

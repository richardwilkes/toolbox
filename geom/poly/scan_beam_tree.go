// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

type scanBeamTree struct {
	root    *scanBeamNode
	entries int
}

type scanBeamNode struct {
	less *scanBeamNode
	more *scanBeamNode
	y    Num
}

func (s *scanBeamTree) add(y Num) {
	s.addToScanBeamTreeAt(&s.root, y)
}

func (s *scanBeamTree) addToScanBeamTreeAt(node **scanBeamNode, y Num) {
	switch {
	case *node == nil:
		*node = &scanBeamNode{y: y}
		s.entries++
	case (*node).y > y:
		s.addToScanBeamTreeAt(&(*node).less, y)
	case (*node).y < y:
		s.addToScanBeamTreeAt(&(*node).more, y)
	default:
	}
}

func (s *scanBeamTree) buildScanBeamTable() []Num {
	table := make([]Num, s.entries)
	if s.root != nil {
		s.root.buildScanBeamTableEntries(0, table)
	}
	return table
}

func (sbt *scanBeamNode) buildScanBeamTableEntries(index int, table []Num) int {
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

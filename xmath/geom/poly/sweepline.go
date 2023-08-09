// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

import "golang.org/x/exp/constraints"

type sweepline[T constraints.Float] []*endpoint[T]

func (s *sweepline[T]) insert(item *endpoint[T]) int {
	if len(*s) == 0 {
		*s = append(*s, item)
		return 0
	}
	*s = append(*s, &endpoint[T]{})
	i := len(*s) - 2
	for i >= 0 && segmentCompare(item, (*s)[i]) {
		(*s)[i+1] = (*s)[i]
		i--
	}
	(*s)[i+1] = item
	return i + 1
}

func (s *sweepline[T]) remove(key *endpoint[T]) {
	for i, el := range *s {
		if el == key {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return
		}
	}
}

func segmentCompare[T constraints.Float](e1, e2 *endpoint[T]) bool {
	switch {
	case e1 == e2:
		return false
	case signedArea(e1.pt, e1.other.pt, e2.pt) != 0 || signedArea(e1.pt, e1.other.pt, e2.other.pt) != 0:
		if e1.pt == e2.pt {
			return e1.below(e2.other.pt)
		}
		if endpointCmp(e1, e2) < 0 {
			return e2.above(e1.pt)
		}
		return e1.below(e2.pt)
	case e1.pt == e2.pt:
		return false
	default:
		return endpointCmp(e1, e2) < 0
	}
}

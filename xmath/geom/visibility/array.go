// Copyright ©2019-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package visibility

type array struct {
	data []int
}

func (a *array) size() int {
	return len(a.data)
}

func (a *array) elem(index int) int {
	return a.data[index]
}

func (a *array) set(index, value int) {
	a.data[index] = value
}

func (a *array) pop() int {
	v := a.data[len(a.data)-1]
	a.data = a.data[:len(a.data)-1]
	return v
}

func (a *array) push(v int) {
	a.data = append(a.data, v)
}

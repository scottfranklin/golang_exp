// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import "testing"

type ceTest struct {
	f   func(in []int) (uint32, error)
	arg []int
	val uint32
}

func normalCE(in []int) (ce uint32, err error) {
	return makeCE(in)
}

func expandCE(in []int) (ce uint32, err error) {
	return makeExpandIndex(in[0])
}

func contractCE(in []int) (ce uint32, err error) {
	return makeContractIndex(ctHandle{in[0], in[1]}, in[2])
}

func decompCE(in []int) (ce uint32, err error) {
	return makeDecompose(in[0], in[1])
}

var ceTests = []ceTest{
	{normalCE, []int{0, 0, 0}, 000},
	{normalCE, []int{0, 30, 3}, 0x1E03},
	{normalCE, []int{100, defaultSecondary, 3}, 0x40006403},
	{normalCE, []int{100, 0, 3}, 0xFFFF}, // non-ignorable primary with non-default secondary
	{normalCE, []int{100, 1, 3}, 0xFFFF},
	{normalCE, []int{1 << maxPrimaryBits, defaultSecondary, 0}, 0xFFFF},
	{normalCE, []int{0, 1 << maxSecondaryBits, 0}, 0xFFFF},
	{normalCE, []int{100, defaultSecondary, 1 << maxTertiaryBits}, 0xFFFF},

	{contractCE, []int{0, 0, 0}, 0x80000000},
	{contractCE, []int{1, 1, 1}, 0x80010021},
	{contractCE, []int{1, (1 << maxNBits) - 1, 1}, 0x8001003F},
	{contractCE, []int{(1 << maxTrieIndexBits) - 1, 1, 1}, 0x8001FFE1},
	{contractCE, []int{1, 1, (1 << maxContractOffsetBits) - 1}, 0xBFFF0021},
	{contractCE, []int{1, (1 << maxNBits), 1}, 0xFFFF},
	{contractCE, []int{(1 << maxTrieIndexBits), 1, 1}, 0xFFFF},
	{contractCE, []int{1, (1 << maxContractOffsetBits), 1}, 0xFFFF},

	{expandCE, []int{0}, 0xC0000000},
	{expandCE, []int{5}, 0xC0000005},
	{expandCE, []int{(1 << maxExpandIndexBits) - 1}, 0xDFFFFFFF},
	{expandCE, []int{1 << maxExpandIndexBits}, 0xFFFF},

	{decompCE, []int{0, 0}, 0xE0000000},
	{decompCE, []int{1, 1}, 0xE0000101},
	{decompCE, []int{0x1F, 0x1F}, 0xE0001F1F},
	{decompCE, []int{256, 0x1F}, 0xFFFF},
	{decompCE, []int{0x1F, 256}, 0xFFFF},
}

func TestColElem(t *testing.T) {
	for i, tt := range ceTests {
		in := make([]int, len(tt.arg))
		copy(in, tt.arg)
		ce, err := tt.f(in)
		if tt.val == 0xFFFF {
			if err == nil {
				t.Errorf("%d: expected error for args %x", i, tt.arg)
			}
			continue
		}
		if err != nil {
			t.Errorf("%d: unexpected error: %v", i, err.Error())
		}
		if ce != tt.val {
			t.Errorf("%d: colElem=%X; want %X", i, ce, tt.val)
		}
	}
}

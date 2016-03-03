package main

import (
	"bytes"
	"fmt"
)

type bitset struct {
	words []int64
}

func (b *bitset) Has(val int) bool {
	q := val / 64
	r := (uint)(val % 64)

	if q >= len(b.words) {
		return false
	}
	return (b.words[q] & (1 << r)) != 0
}

func (b *bitset) Set(val int) {
	//q := val / 64
	//r := (uint)(val % 64)
	q := val >> 6
	r := uint(val & 63)

	// enlarge the words to hold more bits
	if q >= len(b.words) {
		tmp := make([]int64, q+1)
		copy(tmp, b.words)
		b.words = tmp
	}
	b.words[q] = b.words[q] | (1 << r)
}

func (b *bitset) String() string {
	var buf bytes.Buffer
	for i := len(b.words) - 1; i >= 0; i-- {
		fmt.Fprintf(&buf, "%064b", b.words[i])
	}
	return buf.String()
}

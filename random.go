package main

import (
	"math/rand"
)

// range specification, note that min <= max
type IntRange struct {
	min, max int
}

// get next random value within the interval including min and max
func (ir *IntRange) NextRandom() int {

	return rand.Intn(ir.max-ir.min+1) + ir.min
}

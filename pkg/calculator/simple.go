package calculator

import "math"

// SimplePackCalculator uses a single pack size for its calculations
type SimplePackCalculator struct {
	PackSize int
}

// Calculate the required number of packs
func (c SimplePackCalculator) Calculate(quantity int) RequiredPacks {
	packs := make(RequiredPacks)
	packs[c.PackSize] = int(math.Ceil(float64(quantity) / float64(c.PackSize)))
	return packs
}

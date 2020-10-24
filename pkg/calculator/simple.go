package calculator

import "math"

// SimplePackCalculator uses a single pack size for its calculations
type SimplePackCalculator struct {
	PackSize int
}

// Calculate the required number of packs
func (c SimplePackCalculator) Calculate(quantity int) RequiredPacks {
	packs := make(RequiredPacks)
	required := math.Ceil(float64(quantity) / float64(c.PackSize))
	packs[c.PackSize] = int(math.Max(required, 0))
	return packs
}

package calculator

// GraphPackCalculator generates a graph of quantity permutations with the available pack sizes
type GraphPackCalculator struct {
	PackSizes []int
}

// Calculate the required number of packs
func (c GraphPackCalculator) Calculate(quantity int) RequiredPacks {
	packs := make(RequiredPacks)
	return packs
}

package calculator

// PackCalculator returns the packs required to satisfy a quantity of items
type PackCalculator interface {
	Calculate(quantity int) RequiredPacks
}

// RequiredPacks is a map of pack sizes and the number required
type RequiredPacks map[int]int

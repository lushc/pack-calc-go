package calculator

import (
	"math"

	"github.com/go-playground/validator"
)

// SimplePackCalculator uses a single pack size for its calculations
type SimplePackCalculator struct {
	PackSize int `validate:"required,gt=0"`
}

// Calculate the required number of packs
func (c SimplePackCalculator) Calculate(quantity int) (RequiredPacks, error) {
	err := validator.New().Struct(c)
	if err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	packs := make(RequiredPacks)
	required := math.Ceil(float64(quantity) / float64(c.PackSize))
	packs[c.PackSize] = int(math.Max(required, 0))
	return packs, nil
}

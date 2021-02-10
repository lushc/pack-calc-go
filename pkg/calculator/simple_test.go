package calculator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimplePackCalculator(t *testing.T) {
	tests := []struct {
		name     string
		packs    PackCalculator
		quantity int
		want     RequiredPacks
	}{
		{"zero quantity", SimplePackCalculator{250}, 0, RequiredPacks{250: 0}},
		{"negative quantity", SimplePackCalculator{250}, -250, RequiredPacks{250: 0}},
		{"single pack", SimplePackCalculator{250}, 250, RequiredPacks{250: 1}},
		{"divisible", SimplePackCalculator{50}, 500, RequiredPacks{50: 10}},
		{"indivisible", SimplePackCalculator{33}, 500, RequiredPacks{33: 16}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.packs.Calculate(tt.quantity)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSimplePackCalculatorInvalidArguments(t *testing.T) {
	tests := []struct {
		name  string
		packs PackCalculator
	}{
		{"no pack size", SimplePackCalculator{}},
		{"zero pack size", SimplePackCalculator{0}},
		{"negative pack size", SimplePackCalculator{-250}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.packs.Calculate(0)

			if err == nil {
				t.Errorf("%+v.Calculate() == %v, expected error", tt.packs, got)
			}
		})
	}
}

func BenchmarkSimplePackCalculator_Calculate(b *testing.B) {
	benches := []struct {
		name     string
		packs    PackCalculator
		quantity int
	}{
		{"zero quantity", SimplePackCalculator{250}, 0},
		{"negative quantity", SimplePackCalculator{250}, -250},
		{"single pack", SimplePackCalculator{250}, 250},
		{"divisible", SimplePackCalculator{50}, 500},
		{"indivisible", SimplePackCalculator{33}, 500},
	}

	for _, bb := range benches {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = bb.packs.Calculate(bb.quantity)
			}
		})
	}
}

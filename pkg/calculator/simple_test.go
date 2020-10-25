package calculator

import (
	"reflect"
	"testing"
)

func TestSimplePackCalculator(t *testing.T) {
	cases := []struct {
		packs    PackCalculator
		quantity int
		want     RequiredPacks
	}{
		// zero quantity
		{SimplePackCalculator{250}, 0, RequiredPacks{250: 0}},
		// negative quantity
		{SimplePackCalculator{250}, -250, RequiredPacks{250: 0}},
		// single pack
		{SimplePackCalculator{250}, 250, RequiredPacks{250: 1}},
		// divisible
		{SimplePackCalculator{50}, 500, RequiredPacks{50: 10}},
		// undivisible
		{SimplePackCalculator{33}, 500, RequiredPacks{33: 16}},
	}

	for _, c := range cases {
		got, _ := c.packs.Calculate(c.quantity)

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%+v.Calculate(%v) == %v, want %v", c.packs, c.quantity, got, c.want)
		}
	}
}

func TestSimplePackCalculatorInvalidArguments(t *testing.T) {
	cases := []struct {
		packs PackCalculator
	}{
		// no pack size
		{SimplePackCalculator{}},
		// zero pack size
		{SimplePackCalculator{0}},
		// negative pack size
		{SimplePackCalculator{-250}},
	}

	for _, c := range cases {
		got, err := c.packs.Calculate(0)

		if err == nil {
			t.Errorf("%+v.Calculate() == %v, expected error", c.packs, got)
		}
	}
}

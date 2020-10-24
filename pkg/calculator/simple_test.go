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
		got := c.packs.Calculate(c.quantity)

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%+v.Calculate(%v) == %v, want %v", c.packs, c.quantity, got, c.want)
		}
	}
}

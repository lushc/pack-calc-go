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
		{SimplePackCalculator{250}, 0, RequiredPacks{250: 0}},
		{SimplePackCalculator{250}, 250, RequiredPacks{250: 1}},
		{SimplePackCalculator{250}, 251, RequiredPacks{250: 2}},
	}

	for _, c := range cases {
		got := c.packs.Calculate(c.quantity)

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%+v.Calculate(%v) == %v, want %v", c.packs, c.quantity, got, c.want)
		}
	}
}

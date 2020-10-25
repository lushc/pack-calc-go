package calculator

import (
	"reflect"
	"testing"
)

func TestGraphPackCalculator(t *testing.T) {
	cases := []struct {
		packs    PackCalculator
		quantity int
		want     RequiredPacks
	}{
		// default packs with 1
		{GraphPackCalculator{[]int{250, 500, 1000, 2000, 5000}}, 1, RequiredPacks{250: 1}},
		// default packs with 250
		{GraphPackCalculator{[]int{250, 500, 1000, 2000, 5000}}, 250, RequiredPacks{250: 1}},
		// default packs with 251
		{GraphPackCalculator{[]int{250, 500, 1000, 2000, 5000}}, 251, RequiredPacks{500: 1}},
		// default packs with 501
		{GraphPackCalculator{[]int{250, 500, 1000, 2000, 5000}}, 501, RequiredPacks{250: 1, 500: 1}},
		// default packs with 12001
		{GraphPackCalculator{[]int{250, 500, 1000, 2000, 5000}}, 12001, RequiredPacks{250: 1, 2000: 1, 5000: 2}},
		// prime packs with 32
		{GraphPackCalculator{[]int{23, 31, 53, 151, 757}}, 32, RequiredPacks{23: 2}},
		// prime packs with 500
		{GraphPackCalculator{[]int{23, 31, 53, 151, 757}}, 500, RequiredPacks{23: 4, 53: 2, 151: 2}},
		// prime packs with 758
		{GraphPackCalculator{[]int{23, 31, 53, 151, 757}}, 758, RequiredPacks{23: 4, 31: 2, 151: 4}},
		// off by one pack with 500
		{GraphPackCalculator{[]int{1, 100, 200, 499}}, 500, RequiredPacks{1: 1, 499: 1}},
		// edge case pack permutation
		{GraphPackCalculator{[]int{200, 300, 1000}}, 3100, RequiredPacks{200: 1, 300: 3, 1000: 2}},
		// choose smallest pack count
		{GraphPackCalculator{[]int{3, 23, 31, 53, 151, 757}}, 508, RequiredPacks{3: 3, 23: 2, 151: 3}},
		// prime stress test
		{GraphPackCalculator{[]int{23, 31, 53, 151, 757}}, 500000, RequiredPacks{23: 4, 31: 1, 53: 2, 151: 1, 757: 660}},
		// prime stress test with 3 sizes
		{GraphPackCalculator{[]int{23, 31, 53}}, 500000, RequiredPacks{23: 2, 31: 7, 53: 9429}},
		// prime stress test with 2 sizes
		{GraphPackCalculator{[]int{23, 31}}, 500000, RequiredPacks{23: 27, 31: 16109}},
	}

	for _, c := range cases {
		got, _ := c.packs.Calculate(c.quantity)

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%+v.Calculate(%v) == %v, want %v", c.packs, c.quantity, got, c.want)
		}
	}
}

func TestGraphPackCalculatorInvalidArguments(t *testing.T) {
	cases := []struct {
		packs PackCalculator
	}{
		// no pack sizes
		{GraphPackCalculator{[]int{}}},
		// zero pack size
		{GraphPackCalculator{[]int{0}}},
		// negative pack size
		{GraphPackCalculator{[]int{-250}}},
	}

	for _, c := range cases {
		got, err := c.packs.Calculate(0)

		if err == nil {
			t.Errorf("%+v.Calculate() == %v, expected error", c.packs, got)
		}
	}
}

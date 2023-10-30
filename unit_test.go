package main

import (
	"testing"
)

func TestGetNewUnitStrength(t *testing.T) {
	type test struct {
		name string
		unit UnitType
		want int
	}
	tests := []test{
		{name: "Tank", unit: Tank, want: 2},
		{name: "Fighter", unit: Fighter, want: 1},
		{name: "Bomber", unit: Bomber, want: 1},
		{name: "Transport", unit: Transport, want: 3},
		{name: "Destroyer", unit: Destroyer, want: 3},
		{name: "Submarine", unit: Submarine, want: 3},
		{name: "Carrier", unit: Carrier, want: 12},
		{name: "Battleship", unit: Battleship, want: 18},
	}
	//c := NewCity(0, 0, 0)
	for _, tc := range tests {
		got := GetNewUnitStrength(tc.unit)
		//if !reflect.DeepEqual(err, tc.want) {
		if got != tc.want {
			t.Fatalf("GetNewUnitStrength(), name:%s, expected: %d, got: %d", tc.name, tc.want, got)
		}
	}
}

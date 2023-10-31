package main

import (
	"testing"
)

func TestNewCity(t *testing.T) {
	city := NewCity(0, 0)

	if city.OccupyingPlayer != Unoccupied {
		t.Errorf("NewCity() value of OccupyingPlayer = %d; want %d", city.OccupyingPlayer, Unoccupied)
	}
	if city.ManufacturingUnit != Blank {
		t.Errorf("NewCity() value of ManufacturingUnit = %d; want %d", city.ManufacturingUnit, Blank)
	}
}

func TestManufactureUnit(t *testing.T) {
	city := NewCity(0, 0)
	city.OccupyCity(1) // Occupied by player 1
	city.SetManufacturingUnit(Tank)
	day := 0
	if city.ManufacturingUnit != Tank {
		t.Errorf("day %d, occupied city, value of ManufacturingUnit = %d; want %d", day, city.ManufacturingUnit, Tank)
	}
	if city.DaysUntilUnitReady != 4 {
		t.Errorf("day %d, occupied city, value of DaysUntilUnitReady = %d; want %d", day, city.DaysUntilUnitReady, 4)
	}

	day++ // day 1
	unitReady := city.ManufactureUnit()
	if unitReady {
		t.Errorf("day %d, occupied city, value of ready = %t; want %t", day, unitReady, false)
	}
	if city.ManufacturingUnit != Tank {
		t.Errorf("day %d, occupied city, value of ManufacturingUnit = %d; want %d", day, city.ManufacturingUnit, Tank)
	}
	if city.DaysUntilUnitReady != 3 {
		t.Errorf("day %d, occupied city, value of DaysUntilUnitReady = %d; want %d", day, city.DaysUntilUnitReady, 3)
	}

	day++ // day 2
	unitReady = city.ManufactureUnit()
	if unitReady {
		t.Errorf("day %d, occupied city, value of ready = %t; want %t", day, unitReady, false)
	}
	if city.ManufacturingUnit != Tank {
		t.Errorf("day %d, occupied city, value of ManufacturingUnit = %d; want %d", day, city.ManufacturingUnit, Tank)
	}
	if city.DaysUntilUnitReady != 2 {
		t.Errorf("day %d, occupied city, value of DaysUntilUnitReady = %d; want %d", day, city.DaysUntilUnitReady, 2)
	}

	day++ // day 3
	unitReady = city.ManufactureUnit()
	if unitReady {
		t.Errorf("day %d, occupied city, value of ready = %t; want %t", day, unitReady, false)
	}
	if city.ManufacturingUnit != Tank {
		t.Errorf("day %d, occupied city, value of ManufacturingUnit = %d; want %d", day, city.ManufacturingUnit, Tank)
	}
	if city.DaysUntilUnitReady != 1 {
		t.Errorf("day %d, occupied city, value of DaysUntilUnitReady = %d; want %d", day, city.DaysUntilUnitReady, 1)
	}

	day++ // day 4
	unitReady = city.ManufactureUnit()
	if !unitReady {
		t.Errorf("day %d, occupied city, value of ready = %t; want %t", day, unitReady, true)
	}
	if city.ManufacturingUnit != Tank {
		t.Errorf("day %d, occupied city, value of ManufacturingUnit = %d; want %d", day, city.ManufacturingUnit, Tank)
	}
	if city.DaysUntilUnitReady != 4 { // value of city.DaysUntilUnitReady has changed to value of c.GetDaysToProduceUnit(c.ManufacturingUnit)
		t.Errorf("day %d, occupied city, value of DaysUntilUnitReady = %d; want %d", day, city.DaysUntilUnitReady, 4)
	}
}

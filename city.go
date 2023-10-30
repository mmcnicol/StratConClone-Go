package main

// CityState represents the state of a city.
type CityState int

const (
	Unoccupied CityState = iota
	OccupiedByPlayer1
	OccupiedByPlayer2
)

var NewCityStrength int = 2

// City struct represents a city in the game.
type City struct {
	PositionX          int
	PositionY          int
	Strength           int
	OccupyingPlayer    CityState
	ManufacturingUnit  UnitType
	DaysUntilUnitReady int
	IsCityNextToSea    bool
}

// NewCity creates a new City with the given parameters.
func NewCity(positionX, positionY int) *City {
	return &City{
		PositionX:         positionX,
		PositionY:         positionY,
		Strength:          NewCityStrength,
		OccupyingPlayer:   Unoccupied,
		ManufacturingUnit: Blank,
	}
}

// OccupyCity occupies the city by a player.
func (c *City) OccupyCity(player int) {
	c.OccupyingPlayer = CityState(player)
	c.ManufacturingUnit = Blank
	c.DaysUntilUnitReady = 0
}

// SetManufacturingUnit sets the unitTye that the city should manufacture.
func (c *City) SetManufacturingUnit(unit UnitType) {
	c.ManufacturingUnit = unit
	c.DaysUntilUnitReady = c.GetDaysToProduceUnit(unit)
}

// ManufactureUnit updates the days until the unit is ready and returns true if the unit is ready.
func (c *City) ManufactureUnit() bool {
	if c.OccupyingPlayer == Unoccupied {
		return false
	}
	if c.ManufacturingUnit == Blank { // this check should not be required, since an occupied city has to manufacture a unit
		return false
	}
	if c.DaysUntilUnitReady > 0 {
		c.DaysUntilUnitReady--
	}
	if c.DaysUntilUnitReady == 0 {
		c.DaysUntilUnitReady = c.GetDaysToProduceUnit(c.ManufacturingUnit)
		return true
	}
	return false
}

// GetDaysToProduceUnit gets the number of days to produce a unit.
func (c *City) GetDaysToProduceUnit(unit UnitType) int {
	switch unit {
	case Tank:
		return 4
	case Fighter:
		return 6
	case Bomber:
		return 25
	case Transport:
		return 8
	case Destroyer:
		return 8
	case Submarine:
		return 8
	case Carrier:
		return 10
	case Battleship:
		return 20
	}
	return 0
}

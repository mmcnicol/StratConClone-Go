package main

// UnitType represents the type of units that can be manufactured in a city.
type UnitType int

const (
	Blank UnitType = iota
	Tank
	Fighter
	Bomber
	Transport
	Destroyer
	Submarine
	Carrier
	Battleship
)

// Unit struct represents a game unit in the game.
type Unit struct {
	PositionX          int
	PositionY          int
	Type               UnitType
	Player             int
	Strength           int
	MovesPerDay        int
	MovesLeftThisDay   int
	Fuel               int
	CanMoveOnLand      bool
	CanMoveOnWater     bool
	CanFly             bool
	SonarRange         int
	AttackRange        int
	AttacksLeftThisDay int
	CanCaptureCity     bool
}

func NewUnit(positionX, positionY int, unitType UnitType, player int) *Unit {
	return &Unit{
		PositionX:          positionX,
		PositionY:          positionY,
		Type:               unitType,
		Player:             player,
		Strength:           GetNewUnitStrength(unitType),
		MovesPerDay:        2,     // TODO: create a function to get this based on unitType
		MovesLeftThisDay:   2,     // TODO: create a function to get this based on unitType (same as above)
		Fuel:               0,     // TODO: create a function to get this based on unitType
		CanMoveOnLand:      true,  // TODO: create a function to get this based on unitType
		CanMoveOnWater:     false, // TODO: create a function to get this based on unitType
		CanFly:             false, // TODO: create a function to get this based on unitType
		AttackRange:        1,     // TODO: create a function to get this based on unitType
		AttacksLeftThisDay: 2,
		CanCaptureCity:     true, // TODO: create a function to get this based on unitType
	}
}

// GetNewUnitStrength gets the strength of a new unit.
func GetNewUnitStrength(unit UnitType) int {
	switch unit {
	case Tank:
		return 2
	case Fighter:
		return 1
	case Bomber:
		return 1
	case Transport:
		return 3
	case Destroyer:
		return 3
	case Submarine:
		return 3
	case Carrier:
		return 12
	case Battleship:
		return 18
	}
	return 0
}

// MoveTo updates the unit's position on the board, reduces MovesLeftThisDay, and if applicable, reduces Fuel
func (u *Unit) MoveTo(Coordinate Coordinate) {
	u.PositionX = Coordinate.PositionX
	u.PositionY = Coordinate.PositionY
	u.MovesLeftThisDay--
	if u.CanFly {
		u.Fuel--
	}
}

package main

import "fmt"

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

// GetDaysToProduceUnit gets the number of days to produce a unit.
func GetDaysToProduceUnit(unit UnitType) int {
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

// MoveTo updates the unit's position on the board, reduces MovesLeftThisDay, and if applicable, reduces Fuel
func (u *Unit) MoveTo(coordinate Coordinate) {
	fmt.Printf("MoveTo %d, %d\n", coordinate.PositionX, coordinate.PositionY)
	u.PositionX = coordinate.PositionX
	u.PositionY = coordinate.PositionY
	u.MovesLeftThisDay--
	if u.CanFly {
		u.Fuel--
	}
}

// Symbol returns a character depending on the unit type
func (u *Unit) Symbol() string {
	switch u.Type {
	case Tank:
		return "T"
	case Fighter:
		return "F"
	case Bomber:
		return "B"
	case Transport:
		return "R"
	case Destroyer:
		return "D"
	case Submarine:
		return "S"
	case Carrier:
		return "C"
	case Battleship:
		return "L"
	default:
		return "?"
	}
}

func unitTypeToString(unitType UnitType) string {
	switch unitType {
	case Blank:
		return "Blank"
	case Tank:
		return "Tank"
	case Fighter:
		return "Fighter"
	case Bomber:
		return "Bomber"
	case Transport:
		return "Transport"
	case Destroyer:
		return "Destroyer"
	case Submarine:
		return "Submarine"
	case Carrier:
		return "Carrier"
	case Battleship:
		return "Battleship"
	default:
		return "Unknown"
	}
}

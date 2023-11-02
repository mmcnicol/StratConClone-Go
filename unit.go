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
		MovesLeftThisDay:   GetMovesPerDay(unitType),
		Fuel:               GetFuelPerDay(unitType),
		CanMoveOnLand:      GetCanMoveOnLand(unitType),
		CanMoveOnWater:     GetCanMoveOnWater(unitType),
		CanFly:             GetCanFly(unitType),
		AttackRange:        GetAttackRange(unitType),
		AttacksLeftThisDay: GetAttacksPerDay(unitType),
		CanCaptureCity:     GetCanCaptureCity(unitType),
	}
}

// GetCanCaptureCity returns whether or not the unit type can capture a city.
func GetCanCaptureCity(unitType UnitType) bool {
	switch unitType {
	case Tank:
		return true
	case Fighter, Bomber, Transport, Destroyer, Submarine, Carrier, Battleship:
		return false
	default:
		return false
	}
}

// GetAttacksPerDay returns the attack range of the unit type.
func GetAttacksPerDay(unitType UnitType) int {
	switch unitType {
	case Tank, Fighter, Bomber, Transport, Destroyer, Submarine, Carrier, Battleship:
		return 2
	default:
		return 2
	}
}

// GetAttackRange returns the attack range of the unit type.
func GetAttackRange(unitType UnitType) int {
	switch unitType {
	case Tank, Fighter, Bomber, Transport, Destroyer, Submarine, Carrier:
		return 1
	case Battleship:
		return 4
	default:
		return 1
	}
}

// GetCanFly returns whether of not the unit type can fly.
func GetCanFly(unitType UnitType) bool {
	switch unitType {
	case Tank, Transport, Destroyer, Submarine, Carrier, Battleship:
		return false
	case Fighter, Bomber:
		return true
	default:
		return false
	}
}

// GetCanMoveOnWater returns whether of not the unit type can move on water.
func GetCanMoveOnWater(unitType UnitType) bool {
	switch unitType {
	case Tank:
		return false
	case Fighter:
		return true
	case Bomber:
		return true
	case Transport:
		return true
	case Destroyer:
		return true
	case Submarine:
		return true
	case Carrier:
		return true
	case Battleship:
		return true
	default:
		return false
	}
}

// GetCanMoveOnLand returns whether of not the unit type can move on land.
func GetCanMoveOnLand(unitType UnitType) bool {
	switch unitType {
	case Tank:
		return true
	case Fighter:
		return true
	case Bomber:
		return true
	case Transport:
		return false
	case Destroyer:
		return false
	case Submarine:
		return false
	case Carrier:
		return false
	case Battleship:
		return false
	default:
		return false
	}
}

// GetMovesPerDay returns the number of moves per day for a given unit type.
func GetMovesPerDay(unitType UnitType) int {
	switch unitType {
	case Tank:
		return 2
	case Fighter:
		return 20
	case Bomber:
		return 10
	case Transport, Submarine, Carrier, Battleship:
		return 3
	case Destroyer:
		return 4
	}
	return 0
}

// GetFuelPerDay returns the amount of fuel per day for a given unit type.
func GetFuelPerDay(unitType UnitType) int {
	switch unitType {
	case Fighter:
		return 20
	case Bomber:
		return 30
	default:
		return 0
	}
}

// GetNewUnitStrength gets the strength of a new unit.
func GetNewUnitStrength(unitType UnitType) int {
	switch unitType {
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
func GetDaysToProduceUnit(unitType UnitType) int {
	switch unitType {
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
	//fmt.Printf("MoveTo %d, %d\n", coordinate.PositionX, coordinate.PositionY)
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

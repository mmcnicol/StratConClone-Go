package main

import (
	"container/list"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// GameBoard struct represents the game board/grid.
type GameBoard struct {
	Rows    int
	Columns int
	Grid    [][]Cell // 2D slice representing the grid
	Cities  []City
	Units   []Unit
	Day     int
	Player1 *Player
	Player2 *Player
}

// Cell struct represents a cell on the game board.
type Cell struct {
	IsLand  bool // true for land, false for sea
	IsFog   bool // true for fog of war, false if visible
	HasCity bool // true if the cell has a city, false otherwise
}

// Coordinate struct represents an X, Y, position on the game board.
type Coordinate struct {
	PositionX int
	PositionY int
}

type unitWeight struct {
	unit   UnitType
	weight int
}

// NewGameBoard creates a new game board with the specified number of rows and columns.
func NewGameBoard(rows, columns int) *GameBoard {
	grid := make([][]Cell, rows)
	for i := range grid {
		row := make([]Cell, columns)
		for j := range row {
			row[j].IsFog = true // Initialize Fog to true for all cells
			// the default value of HasCity is false by default, which is appropriate
			// the default value of IsLand is 0 (false) by default, which represents Sea
		}
		grid[i] = row
	}

	return &GameBoard{
		Rows:    rows,
		Columns: columns,
		Grid:    grid,
		Day:     0,
	}
}

// GenerateRandomIslands generates random oval-shaped islands on the game board.
func (g *GameBoard) GenerateRandomIslands(numIslands int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < numIslands; i++ {
		centerRow := r.Intn(g.Rows)
		centerCol := r.Intn(g.Columns)
		radiusX := r.Intn(8) + 2 // Random oval radius between 2 and 5 cells
		radiusY := r.Intn(6) + 2 // Random oval radius between 2 and 5 cells

		for row := 0; row < g.Rows; row++ {
			for col := 0; col < g.Columns; col++ {
				dx := float64(col - centerCol)
				dy := float64(row - centerRow)
				distance := math.Pow(dx/float64(radiusX), 2) + math.Pow(dy/float64(radiusY), 2)

				if distance <= 1.0 {
					g.Grid[row][col].IsLand = true
				}
			}
		}
	}
}

// AddCities randomly adds cities to land cells without neighboring cities.
func (g *GameBoard) AddCities(numCities int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < numCities; i++ {
		for {
			row := r.Intn(g.Rows)
			col := r.Intn(g.Columns)

			// Check if the cell is land and does not have neighboring cities
			excludeTargetCell := false
			if g.Grid[row][col].IsLand && !g.HasNeighboringCity(row, col, excludeTargetCell) {
				g.Grid[row][col].HasCity = true
				city := NewCity(row, col)
				city.IsCityNextToSea = g.IsCityNextToSea(city.PositionX, city.PositionY)
				g.Cities = append(g.Cities, *city)
				break
			}
		}
	}
}

// DayZero performs game logic for a new day
func (g *GameBoard) DayZero() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// player 1
	randomIndex := r.Intn(len(g.Cities))
	city := &g.Cities[randomIndex]
	//city.OccupyingPlayer = OccupiedByPlayer1
	city.OccupyCity(1)
	if g.Player1.IsAI {
		city.SetManufacturingUnit(g.getWhichUnitToManufactureNextAI(Coordinate{city.PositionX, city.PositionY}, 1, city.IsCityNextToSea))
	}
	// player 2
	randomIndex = r.Intn(len(g.Cities))
	city = &g.Cities[randomIndex]
	//city.OccupyingPlayer = OccupiedByPlayer2
	city.OccupyCity(2)
	if g.Player2.IsAI {
		city.SetManufacturingUnit(g.getWhichUnitToManufactureNextAI(Coordinate{city.PositionX, city.PositionY}, 2, city.IsCityNextToSea))
	}
}

// NextDay performs game logic for a new day
func (g *GameBoard) NextDay() {
	g.Day++
	for i := range g.Units {
		unit := &g.Units[i] // Get a pointer to the current unit
		unit.MovesLeftThisDay = GetMovesPerDay(unit.Type)
	}
	for i := range g.Cities {
		city := &g.Cities[i] // Get a pointer to the current city
		unitReady := city.ManufactureUnit()
		if unitReady {
			player := 1
			if city.OccupyingPlayer == OccupiedByPlayer2 {
				player = 2
			}
			newUnit := NewUnit(city.PositionX, city.PositionY, city.ManufacturingUnit, player)
			g.Units = append(g.Units, *newUnit)
			// Reset DaysUntilUnitReady to the production time when the unit is manufactured
			city.DaysUntilUnitReady = GetDaysToProduceUnit(city.ManufacturingUnit)
		}
	}
}

// hasNeighboringCity checks if a cell has neighboring cities.
func (g *GameBoard) HasNeighboringCity(row, col int, excludeTargetCell bool) bool {
	for i := row - 1; i <= row+1; i++ {
		for j := col - 1; j <= col+1; j++ {
			// Skip the current cell represented by the arguments
			if excludeTargetCell && i == row && j == col {
				continue
			}
			if i >= 0 && i < g.Rows && j >= 0 && j < g.Columns && g.Grid[i][j].HasCity {
				return true
			}
		}
	}
	return false
}

// isCityNextToSea checks if a city is next to the sea.
func (g *GameBoard) IsCityNextToSea(row, col int) bool {
	for i := row - 1; i <= row+1; i++ {
		for j := col - 1; j <= col+1; j++ {
			if i >= 0 && i < g.Rows && j >= 0 && j < g.Columns && !g.Grid[i][j].HasCity && !g.Grid[i][j].IsLand {
				return true
			}
		}
	}
	return false
}

// IterateGrid iterates over all cells in the grid and applies the given callback function.
func (g *GameBoard) IterateGrid(callback func(row, col int, cell *Cell)) {
	for i := 0; i < g.Rows; i++ {
		for j := 0; j < g.Columns; j++ {
			callback(i, j, &g.Grid[i][j])
		}
	}
}

// Print prints the game board.
func (g *GameBoard) Print(showFogOfWar bool) {
	// Print the game board with land/sea and fog of war
	if showFogOfWar {
		fmt.Println("Game Board (Land/Sea, Cities, Fog of War):")
	} else {
		fmt.Println("Game Board (Land/Sea, Cities):")
	}
	for i := 0; i < g.Rows; i++ {
		for j := 0; j < g.Columns; j++ {
			if showFogOfWar && g.Grid[i][j].IsFog {
				fmt.Print("? ")
			} else {
				if g.Grid[i][j].HasCity {
					fmt.Print("C ")
				} else if g.Grid[i][j].IsLand {
					fmt.Print("L ")
				} else {
					fmt.Print("S ")
				}
			}
		}
		fmt.Println()
	}
}

func (g *GameBoard) printGridWithUnits(showFogOfWar bool) {
	grid := g.printToSlice(showFogOfWar)
	for _, unit := range g.Units {
		if !showFogOfWar || !g.Grid[unit.PositionX][unit.PositionY].IsFog {
			grid[unit.PositionX][unit.PositionY] = unit.Symbol()
		}
	}
	g.printSlice(grid)
}

func (g *GameBoard) printSlice(grid [][]string) {
	for i := 0; i < g.Rows; i++ {
		for j := 0; j < g.Columns; j++ {
			//fmt.Print(grid[i][j] + " ")
			fmt.Print(grid[i][j])
		}
		fmt.Println()
	}
}

func (g *GameBoard) printToSlice(showFogOfWar bool) [][]string {
	grid := make([][]string, g.Rows)
	for i := range grid {
		row := make([]string, g.Columns)
		grid[i] = row
	}
	// Print the game board with land/sea and fog of war
	for i := 0; i < g.Rows; i++ {
		for j := 0; j < g.Columns; j++ {
			if showFogOfWar && g.Grid[i][j].IsFog {
				grid[i][j] = "?"
			} else {
				if g.Grid[i][j].HasCity {
					grid[i][j] = "C"
				} else if g.Grid[i][j].IsLand {
					grid[i][j] = "L"
				} else {
					grid[i][j] = "S"
				}
			}
		}
	}
	return grid
}

func (g *GameBoard) printCitiesForPlayer(playerID int) {
	//fmt.Printf("Cities for Player %d:\n", playerID)
	for _, city := range g.Cities {
		if city.OccupyingPlayer == OccupiedByPlayer1 && playerID == 1 {
			manufacturingUnit := unitTypeToString(city.ManufacturingUnit)
			fmt.Printf("City at (%d, %d) is manufacturing: %s, DaysUntilUnitReady: %d\n", city.PositionX, city.PositionY, manufacturingUnit, city.DaysUntilUnitReady)
		}
		if city.OccupyingPlayer == OccupiedByPlayer2 && playerID == 2 {
			manufacturingUnit := unitTypeToString(city.ManufacturingUnit)
			fmt.Printf("City at (%d, %d) is manufacturing: %s, DaysUntilUnitReady: %d\n", city.PositionX, city.PositionY, manufacturingUnit, city.DaysUntilUnitReady)
		}
	}
}

func (g *GameBoard) DoPlayerTurnAI(player int) {
	var activeUnit *Unit
	//showFogOfWar := true
	//coordinate := Coordinate{}
	for {
		activeUnit = g.getActiveUnitForPlayer(player)
		if activeUnit == nil {
			break // No more active units for the player
		}
		//coordinate = Coordinate{activeUnit.PositionX, activeUnit.PositionY}

		g.runUnitAI(activeUnit)

		if g.hasPlayerWon(player) {
			fmt.Printf("\nDay: %d\n", g.Day)
			fmt.Printf("\nPlayer %d has won\n", player)
			break // the player has won
		}
	}
	/*
		islandMap := g.getIslandMap(coordinate)
		isConquered := g.isIslandConquered(islandMap, player)
		tankCount := g.getUnitCount(Tank, islandMap, player)
		transportCount := g.getUnitCount(Transport, islandMap, player)

		fmt.Printf("\nDay %d, Player %d:, hasConqueredIsland:%t \n", g.Day, player, isConquered)
		fmt.Printf("\nTanks:%d, Transports:%d:\n", tankCount, transportCount)
		g.printCitiesForPlayer(player)
		g.printGridWithUnits(showFogOfWar)
		time.Sleep(20 * time.Millisecond)
		//clearScreen()
	*/
}

// getActiveUnitForPlayer returns an active unit for the specified player with MovesLeftThisDay > 0.
func (g *GameBoard) getActiveUnitForPlayer(player int) *Unit {
	for i := range g.Units {
		if g.Units[i].Player == player && g.Units[i].MovesLeftThisDay > 0 {
			return &g.Units[i]
		}
	}
	return nil
}

// runUnitAI implements the AI logic for the a unit.
func (g *GameBoard) runUnitAI(unit *Unit) {
	possibleMoves := g.getPossibleMoves(unit)
	if len(possibleMoves) > 0 {
		move := possibleMoves[rand.Intn(len(possibleMoves))]
		g.attemptMoveTo(move, unit)
	}
}

/*
// getPossibleMoves returns possible moves for the given unit.
//
// priorities:
// should attack enemy unit if nearby
// should attack enemy city if nearby
// should attack unoccupied city if nearby
// should move to clear fog of war if nearby
// if island has been conquered, move to staging point and wait
func (g *GameBoard) getPossibleMoves(unit *Unit) []Coordinate {
	var moves []Coordinate
	var enemyUnits []Coordinate
	var enemyCities []Coordinate
	var unoccupiedCities []Coordinate
	var fogOfWar []Coordinate
	var randomMoves []Coordinate

	// Check neighboring cells and add valid moves to the list
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue // Skip the current cell
			}
			newRow, newCol := unit.PositionX+i, unit.PositionY+j
			if newRow >= 0 && newRow < g.Rows && newCol >= 0 && newCol < g.Columns {
				defender := g.getUnitAtCoordinates(Coordinate{newRow, newCol}, unit.Player)
				if defender != nil {
					enemyUnits = append(enemyUnits, Coordinate{newRow, newCol})
				}
				if g.Grid[newRow][newCol].HasCity {
					city := g.getCityAtCoordinates(Coordinate{newRow, newCol})
					if city.OccupyingPlayer == Unoccupied {
						unoccupiedCities = append(unoccupiedCities, Coordinate{newRow, newCol})
					} else if city.OccupyingPlayer == OccupiedByPlayer1 && unit.Player != 1 {
						enemyCities = append(enemyCities, Coordinate{newRow, newCol})
					} else if city.OccupyingPlayer == OccupiedByPlayer2 && unit.Player != 2 {
						enemyCities = append(enemyCities, Coordinate{newRow, newCol})
					}
				}
				if g.Grid[newRow][newCol].IsFog {
					fogOfWar = append(fogOfWar, Coordinate{newRow, newCol})
				}
				randomMoves = append(randomMoves, Coordinate{newRow, newCol})
			}
		}
	}
	if len(enemyUnits) > 0 {
		moves = append(moves, enemyUnits[0])
	} else if len(enemyCities) > 0 {
		moves = append(moves, enemyCities[0])
	} else if len(unoccupiedCities) > 0 {
		moves = append(moves, unoccupiedCities[0])
	} else if len(fogOfWar) > 0 {
		moves = append(moves, fogOfWar[0])
	} else {
		islandMap := g.getIslandMap(Coordinate{unit.PositionX, unit.PositionY})
		isConquered := g.isIslandConquered(islandMap, unit.Player)
		if isConquered {
			stagingPoint := g.getIsIslandCityNextToSea(islandMap)
			if stagingPoint != nil {
				pathToStagingPoint := g.FindPath(*stagingPoint, unit)
				if pathToStagingPoint != nil {
					firstStepOnPathTowardsStagingPoint := getSecondCoordinate(pathToStagingPoint)
					if firstStepOnPathTowardsStagingPoint != nil {
						moves = append(moves, *firstStepOnPathTowardsStagingPoint)
					}
				}
			}
		}
	}
	if len(moves) == 0 {
		moves = append(moves, randomMoves...)
	}
	return moves
}
*/

// getPossibleMoves returns possible moves for the given unit.
func (g *GameBoard) getPossibleMoves(unit *Unit) []Coordinate {
	var moves []Coordinate

	enemyUnits := g.getEnemyUnitsCoordinates(unit)
	enemyCities := g.getEnemyCitiesCoordinates(unit)
	unoccupiedCities := g.getUnoccupiedCitiesCoordinates(unit)
	fogOfWar := g.getFogOfWarCoordinates(unit)
	stagingPoint := g.getStagingPoint(unit)
	randomMoves := g.getRandomMoves(unit)

	if len(fogOfWar) > 0 {
		moves = append(moves, fogOfWar[0])
	} else if len(enemyUnits) > 0 {
		moves = append(moves, enemyUnits[0])
	} else if len(enemyCities) > 0 {
		moves = append(moves, enemyCities[0])
	} else if len(unoccupiedCities) > 0 {
		moves = append(moves, unoccupiedCities[0])
	} else if len(stagingPoint) > 0 {
		moves = append(moves, stagingPoint[0])
	} else {
		moves = append(moves, randomMoves...)
	}

	return moves
}

func (g *GameBoard) getEnemyUnitsCoordinates(unit *Unit) []Coordinate {
	var moves []Coordinate
	if unit.CanFly || unit.CanMoveOnWater {
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				if i == 0 && j == 0 {
					continue // Skip the current cell
				}
				newRow, newCol := unit.PositionX+i, unit.PositionY+j
				if newRow >= 0 && newRow < g.Rows && newCol >= 0 && newCol < g.Columns {
					defender := g.getUnitAtCoordinates(Coordinate{newRow, newCol}, unit.Player)
					if defender != nil {
						moves = append(moves, Coordinate{newRow, newCol})
					}
				}
			}
		}
	} else if unit.CanMoveOnLand {
		islandMap := g.getIslandMap(Coordinate{unit.PositionX, unit.PositionY})
		isConquered := g.isIslandConquered(islandMap, unit.Player)
		if !isConquered {
			destination := g.getIsIslandEnemyUnit(islandMap, unit)
			if destination != nil {
				pathToDestination := g.FindPath(*destination, unit)
				if pathToDestination != nil {
					firstStepOnPathTowardsDestination := getSecondCoordinate(pathToDestination)
					if firstStepOnPathTowardsDestination != nil {
						moves = append(moves, *firstStepOnPathTowardsDestination)
					}
				}
			}
		}
	}
	return moves
}

func (g *GameBoard) getEnemyCitiesCoordinates(unit *Unit) []Coordinate {
	var moves []Coordinate
	// Logic to find enemy cities
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue // Skip the current cell
			}
			newRow, newCol := unit.PositionX+i, unit.PositionY+j
			if newRow >= 0 && newRow < g.Rows && newCol >= 0 && newCol < g.Columns {
				if g.Grid[newRow][newCol].HasCity {
					city := g.getCityAtCoordinates(Coordinate{newRow, newCol})
					if city.OccupyingPlayer == OccupiedByPlayer1 && unit.Player != 1 {
						moves = append(moves, Coordinate{newRow, newCol})
					} else if city.OccupyingPlayer == OccupiedByPlayer2 && unit.Player != 2 {
						moves = append(moves, Coordinate{newRow, newCol})
					}
				}
			}
		}
	}
	return moves
}

func (g *GameBoard) getUnoccupiedCitiesCoordinates(unit *Unit) []Coordinate {
	var moves []Coordinate
	if unit.CanCaptureCity {
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				if i == 0 && j == 0 {
					continue // Skip the current cell
				}
				newRow, newCol := unit.PositionX+i, unit.PositionY+j
				if newRow >= 0 && newRow < g.Rows && newCol >= 0 && newCol < g.Columns {
					if g.Grid[newRow][newCol].HasCity {
						city := g.getCityAtCoordinates(Coordinate{newRow, newCol})
						if city.OccupyingPlayer == Unoccupied {
							moves = append(moves, Coordinate{newRow, newCol})
						}
					}
				}
			}
		}
	}
	return moves
}

func (g *GameBoard) getFogOfWarCoordinates(unit *Unit) []Coordinate {
	var moves []Coordinate
	if unit.CanFly || unit.CanMoveOnWater {
		// Logic to find fog of war cells
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				if i == 0 && j == 0 {
					continue // Skip the current cell
				}
				newRow, newCol := unit.PositionX+i, unit.PositionY+j
				if newRow >= 0 && newRow < g.Rows && newCol >= 0 && newCol < g.Columns {
					if g.Grid[newRow][newCol].IsFog {
						moves = append(moves, Coordinate{newRow, newCol})
					}
				}
			}
		}
	} else if unit.CanMoveOnLand {
		islandMap := g.getIslandMap(Coordinate{unit.PositionX, unit.PositionY})
		isConquered := g.isIslandConquered(islandMap, unit.Player)
		if !isConquered {
			destination := g.getIsIslandFogOfWar(islandMap)
			if destination != nil {
				pathToDestination := g.FindPath(*destination, unit)
				if pathToDestination != nil {
					firstStepOnPathTowardsDestination := getSecondCoordinate(pathToDestination)
					if firstStepOnPathTowardsDestination != nil {
						moves = append(moves, *firstStepOnPathTowardsDestination)
					}
				}
			}
		}
	}
	return moves
}

// getStagingPoint returns coordinate on path towards staging point
func (g *GameBoard) getStagingPoint(unit *Unit) []Coordinate {
	var moves []Coordinate
	if unit.Type == Tank {
		islandMap := g.getIslandMap(Coordinate{unit.PositionX, unit.PositionY})
		isConquered := g.isIslandConquered(islandMap, unit.Player)
		if isConquered {
			stagingPoint := g.getIsIslandCityNextToSea(islandMap)
			if stagingPoint != nil {
				pathToStagingPoint := g.FindPath(*stagingPoint, unit)
				if pathToStagingPoint != nil {
					firstStepOnPathTowardsStagingPoint := getSecondCoordinate(pathToStagingPoint)
					if firstStepOnPathTowardsStagingPoint != nil {
						moves = append(moves, *firstStepOnPathTowardsStagingPoint)
					}
				}
			}
		}
	}
	return moves
}

func (g *GameBoard) getRandomMoves(unit *Unit) []Coordinate {
	var moves []Coordinate
	// Logic to generate random moves
	// Check neighboring cells and add valid moves to the list
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue // Skip the current cell
			}
			newRow, newCol := unit.PositionX+i, unit.PositionY+j
			if newRow >= 0 && newRow < g.Rows && newCol >= 0 && newCol < g.Columns {
				if unit.CanFly {
					moves = append(moves, Coordinate{newRow, newCol})
				} else if g.Grid[newRow][newCol].IsLand && unit.CanMoveOnLand {
					moves = append(moves, Coordinate{newRow, newCol})
				} else if !g.Grid[newRow][newCol].IsLand && unit.CanMoveOnWater {
					moves = append(moves, Coordinate{newRow, newCol})
				}
			}
		}
	}
	return moves
}

// getSecondCoordinate returns the second element in the slice of Coordinate.
func getSecondCoordinate(coordinates []Coordinate) *Coordinate {
	if len(coordinates) > 1 {
		firstCoordinate := coordinates[1]
		return &firstCoordinate
	}
	return nil // Return nil if the slice is empty
}

/*
// attemptMoveTo attempts to move the unit to the destination coordinates
func (g *GameBoard) attemptMoveTo(destinationCoordinate Coordinate, unit *Unit) {
	fmt.Printf("unit at %d, %d, attemptMoveTo() %d, %d\n", unit.PositionX, unit.PositionY, destinationCoordinate.PositionX, destinationCoordinate.PositionY)
	radius := 1
	g.clearFogOfWarAroundCoordinate(destinationCoordinate, radius)
	defender := g.getUnitAtCoordinates(destinationCoordinate, unit.Player)
	if defender != nil {
		g.resolveUnitAttack(unit, defender, g.getAttackOutcome())
	} else if g.Grid[destinationCoordinate.PositionX][destinationCoordinate.PositionY].HasCity && unit.CanMoveOnLand {
		defender := g.getCityAtCoordinates(destinationCoordinate)
		g.resolveCityAttack(unit, defender, g.getAttackOutcome())
	} else if g.Grid[destinationCoordinate.PositionX][destinationCoordinate.PositionY].IsLand && unit.CanMoveOnLand {
		unit.MoveTo(destinationCoordinate)
	} else if !g.Grid[destinationCoordinate.PositionX][destinationCoordinate.PositionY].IsLand && unit.CanMoveOnWater {
		unit.MoveTo(destinationCoordinate)
	} else {
		fmt.Println("attemptMoveTo() illegal move!")
	}
}
*/

// ActionType represents the type of action to be performed
type ActionType int

const (
	// ActionMove represents a move action
	ActionMove ActionType = iota
	// ActionUnitAttack represents a unit attack action
	ActionUnitAttack
	// ActionCityAttack represents a city attack action
	ActionCityAttack
	// ActionIllegalMove represents an illegal move action
	ActionIllegalMove
)

// determineAction determines the action to be performed based on the destination coordinate and unit's properties
func (g *GameBoard) determineAction(destinationCoordinate Coordinate, unit *Unit) ActionType {
	defender := g.getUnitAtCoordinates(destinationCoordinate, unit.Player)
	if defender != nil {
		return ActionUnitAttack
	} else if g.Grid[destinationCoordinate.PositionX][destinationCoordinate.PositionY].HasCity && unit.CanMoveOnLand {
		return ActionCityAttack
	} else if unit.CanFly {
		return ActionMove
		//} else if g.getCityAtCoordinates(destinationCoordinate) !=nil { // any unit can move into a city by water
		//	return ActionMove
	} else if g.Grid[destinationCoordinate.PositionX][destinationCoordinate.PositionY].IsLand && unit.CanCaptureCity {
		return ActionMove
	} else if !g.Grid[destinationCoordinate.PositionX][destinationCoordinate.PositionY].IsLand && unit.CanMoveOnWater {
		return ActionMove
	}
	return ActionIllegalMove
}

// performAction performs the specified action based on the ActionType
func (g *GameBoard) performAction(actionType ActionType, destinationCoordinate Coordinate, unit *Unit) {
	switch actionType {
	case ActionMove:
		unit.MoveTo(destinationCoordinate)
	case ActionUnitAttack:
		defender := g.getUnitAtCoordinates(destinationCoordinate, unit.Player)
		g.resolveUnitAttack(unit, defender, g.getAttackOutcome())
	case ActionCityAttack:
		defender := g.getCityAtCoordinates(destinationCoordinate)
		g.resolveCityAttack(unit, defender, g.getAttackOutcome())
	case ActionIllegalMove:
		fmt.Println("Illegal move!")
	}
}

// attemptMoveTo attempts to move the unit to the destination coordinates
func (g *GameBoard) attemptMoveTo(destinationCoordinate Coordinate, unit *Unit) {
	//fmt.Printf("unit at %d, %d, AttemptMoveTo() %d, %d\n", unit.PositionX, unit.PositionY, destinationCoordinate.PositionX, destinationCoordinate.PositionY)
	radius := 1
	g.clearFogOfWarAroundCoordinate(destinationCoordinate, radius)
	actionType := g.determineAction(destinationCoordinate, unit)
	g.performAction(actionType, destinationCoordinate, unit)
}

// getAttackOutcome decides an attack outcome based on chance
func (g *GameBoard) getAttackOutcome() bool {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(2) == 0 // 50% probability.
}

// clearFogOfWarAroundCoordinate clears the fog of war around the specified coordinates within a given radius.
func (g *GameBoard) clearFogOfWarAroundCoordinate(coordinate Coordinate, radius int) {
	for i := coordinate.PositionX - radius; i <= coordinate.PositionX+radius; i++ {
		for j := coordinate.PositionY - radius; j <= coordinate.PositionY+radius; j++ {
			if i >= 0 && i < g.Rows && j >= 0 && j < g.Columns {
				g.Grid[i][j].IsFog = false
			}
		}
	}
}

// getUnitAtCoordinates retrieves an enemy unit at the specified coordinates.
func (g *GameBoard) getUnitAtCoordinates(coordinate Coordinate, attackingPlayer int) *Unit {
	for i := range g.Units {
		if g.Units[i].PositionX == coordinate.PositionX && g.Units[i].PositionY == coordinate.PositionY && g.Units[i].Player != attackingPlayer {
			return &g.Units[i]
		}
	}
	return nil
}

// getCityAtCoordinates retrieves a city at the specified coordinates.
func (g *GameBoard) getCityAtCoordinates(coordinate Coordinate) *City {
	for i := range g.Cities {
		if g.Cities[i].PositionX == coordinate.PositionX && g.Cities[i].PositionY == coordinate.PositionY {
			return &g.Cities[i]
		}
	}
	return nil
}

// resolveCityAttack determines the outcome of an attack between an attacking unit and a defending city.
func (g *GameBoard) resolveCityAttack(attacker *Unit, defender *City, attackOutcome bool) {
	fmt.Printf("resolveCityAttack defender %d, %d\n", defender.PositionX, defender.PositionY)
	attacker.MovesLeftThisDay--
	if attacker.CanFly {
		attacker.Fuel--
	}
	//if attackOutcome && attacker.Strength >= defender.Strength {
	if attackOutcome {
		// Apply damage to the defender's strength
		defender.Strength--
		// Check if the defender is destroyed
		if defender.Strength <= 0 {
			// Defender is conquered, change OccupyingPlayer
			fmt.Println("Defender [City] is conquered")
			if attacker.Player == 1 {
				defender.OccupyingPlayer = OccupiedByPlayer1
			} else {
				defender.OccupyingPlayer = OccupiedByPlayer2
			}
			defender.Strength = NewCityStrength
			//defender.ManufacturingUnit = Blank // TODO: if player is computer, decide what to manufacture
			defender.ManufacturingUnit = g.getWhichUnitToManufactureNextAI(Coordinate{defender.PositionX, defender.PositionY}, attacker.Player, defender.IsCityNextToSea)
			defender.DaysUntilUnitReady = GetDaysToProduceUnit(defender.ManufacturingUnit)
			// Attacker is destroyed when it conquers a city
			g.removeUnit(attacker)
		}
	} else {
		// Apply damage to the attacker's strength
		attacker.Strength--
		// Check if the attacker is destroyed
		if attacker.Strength <= 0 {
			// Attacker is destroyed, remove it from the game board
			fmt.Println("Attacker is destroyed")
			g.removeUnit(attacker)
		}
	}
}

// resolveUnitAttack determines the outcome of an attack between an attacking unit and a defending unit.
func (g *GameBoard) resolveUnitAttack(attacker, defender *Unit, attackOutcome bool) {
	fmt.Printf("resolveUnitAttack defender %d, %d\n", defender.PositionX, defender.PositionY)
	attacker.MovesLeftThisDay--
	if attacker.CanFly {
		attacker.Fuel--
	}
	if attackOutcome && attacker.Strength >= defender.Strength {
		// Apply damage to the defender's strength
		defender.Strength--
		// Check if the defender is destroyed
		if defender.Strength <= 0 {
			// Defender is destroyed, remove it from the game board
			fmt.Println("Defender is destroyed")
			g.removeUnit(defender)
		}
		// attacker does not move to defenders coordinates
		//attacker.PositionX = defender.PositionX
		//attacker.PositionX = defender.PositionY
	} else {
		// Apply damage to the attacker's strength
		attacker.Strength--
		// Check if the attacker is destroyed
		if attacker.Strength <= 0 {
			// Attacker is destroyed, remove it from the game board
			fmt.Println("Attacker is destroyed")
			g.removeUnit(attacker)
		}
	}
}

// removeUnit removes a unit from the game board's Units slice.
func (g *GameBoard) removeUnit(unitToRemove *Unit) {
	var updatedUnits []Unit
	for _, unit := range g.Units {
		if unit != *unitToRemove {
			updatedUnits = append(updatedUnits, unit)
		}
	}
	g.Units = updatedUnits
}

// getIslandMap returns a slice of coordinates representing the island connected to the given coordinate.
func (g *GameBoard) getIslandMap(coordinate Coordinate) []Coordinate {
	visited := make(map[Coordinate]bool)
	islandMap := make([]Coordinate, 0)

	// Define a recursive flood fill function to explore land cells
	var floodFill func(x, y int)
	floodFill = func(x, y int) {
		// Check if the cell is within the grid boundaries and is a land cell
		if x >= 0 && x < g.Rows && y >= 0 && y < g.Columns && g.Grid[x][y].IsLand {
			// Mark the cell as visited
			visited[Coordinate{PositionX: x, PositionY: y}] = true
			// Add the coordinate to the island map
			islandMap = append(islandMap, Coordinate{PositionX: x, PositionY: y})

			// Explore neighboring cells
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					if i != 0 || j != 0 {
						neighborX, neighborY := x+i, y+j
						neighborCoord := Coordinate{PositionX: neighborX, PositionY: neighborY}
						if !visited[neighborCoord] {
							floodFill(neighborX, neighborY) // Recur for neighboring cell
						}
					}
				}
			}
		}
	}

	// Start the flood fill from the given coordinates
	floodFill(coordinate.PositionX, coordinate.PositionY)

	return islandMap
}

// isIslandConquered checks if all cities on the island represented by coordinates are occupied by the same player.
func (g *GameBoard) isIslandConquered(islandMap []Coordinate, playerID int) bool {
	for _, coord := range islandMap {
		city := g.getCityAtCoordinates(coord)
		if city != nil {
			if city.OccupyingPlayer == Unoccupied {
				return false // Island is not conquered if any city is unoccupied
			} else if city.OccupyingPlayer == OccupiedByPlayer2 && playerID == 1 {
				return false // Island is not conquered if any city is not occupied by the player
			} else if city.OccupyingPlayer == OccupiedByPlayer1 && playerID == 2 {
				return false // Island is not conquered if any city is not occupied by the player
			}
		}
	}
	return true // Island is conquered if all cities are occupied by the player
}

// getIsIslandEnemyUnit returns coordinate of enemy unit on island.
func (g *GameBoard) getIsIslandEnemyUnit(islandMap []Coordinate, attacker *Unit) *Coordinate {
	for _, coord := range islandMap {
		for _, unit := range g.Units {
			if coord.PositionX == unit.PositionX &&
				coord.PositionY == unit.PositionY &&
				attacker.Player != unit.Player {
				return &Coordinate{PositionX: unit.PositionX, PositionY: unit.PositionY}
			}
		}
	}
	return nil // city next to sea not found on island
}

// getIsIslandFogOfWar returns coordinate of fog of war.
func (g *GameBoard) getIsIslandFogOfWar(islandMap []Coordinate) *Coordinate {
	for _, coord := range islandMap {
		if coord.PositionX >= 0 && coord.PositionX < g.Rows &&
			coord.PositionY >= 0 && coord.PositionY < g.Columns &&
			g.Grid[coord.PositionX][coord.PositionY].IsFog {
			return &Coordinate{PositionX: coord.PositionX, PositionY: coord.PositionY}
		}
	}
	return nil // city next to sea not found on island
}

// getIsIslandCityNextToSea returns coordinate of city which is next to sea.
func (g *GameBoard) getIsIslandCityNextToSea(islandMap []Coordinate) *Coordinate {
	for _, coord := range islandMap {
		city := g.getCityAtCoordinates(coord)
		if city != nil && city.IsCityNextToSea {
			return &Coordinate{PositionX: city.PositionX, PositionY: city.PositionY}
		}
	}
	return nil // city next to sea not found on island
}

// FindPath finds a path for the unit to reach the target coordinate on the grid.
func (g *GameBoard) FindPath(target Coordinate, unit *Unit) []Coordinate {
	// Define possible moves: up, down, left, right, ...
	moves := []Coordinate{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {1, 1}, {0, 0}}

	// Initialize visited map to keep track of visited cells
	visited := make(map[Coordinate]bool)
	visited[Coordinate{unit.PositionX, unit.PositionY}] = true

	// Initialize queue for BFS
	queue := list.New()
	queue.PushBack([]Coordinate{{unit.PositionX, unit.PositionY}})

	for queue.Len() > 0 {
		// Dequeue the path from the queue
		path := queue.Remove(queue.Front()).([]Coordinate)
		lastPos := path[len(path)-1]

		// If the last position is the target, return the path
		if lastPos == target {
			return path
		}

		// Explore possible moves
		for _, move := range moves {
			newX, newY := lastPos.PositionX+move.PositionX, lastPos.PositionY+move.PositionY
			newPos := Coordinate{newX, newY}

			// Check if the new position is within the grid boundaries and not visited
			if newX >= 0 && newX < g.Rows && newY >= 0 && newY < g.Columns &&
				!visited[newPos] &&
				(unit.CanFly ||
					(unit.CanMoveOnLand && g.Grid[newX][newY].IsLand) ||
					(unit.CanMoveOnWater && !g.Grid[newX][newY].IsLand)) {

				// Mark the new position as visited
				visited[newPos] = true

				// Enqueue the new path
				newPath := make([]Coordinate, len(path))
				copy(newPath, path)
				newPath = append(newPath, newPos)
				queue.PushBack(newPath)
			}
		}
	}

	// No path found
	return nil
}

// getWhichUnitToManufactureNextAI determine which unit type a city should manufacture next AI
func (g *GameBoard) getWhichUnitToManufactureNextAI(coordinate Coordinate, player int, isCityNextToSea bool) UnitType {
	islandMap := g.getIslandMap(coordinate)
	isConquered := g.isIslandConquered(islandMap, player)
	tankCount := g.getUnitCount(Tank, islandMap, player)
	var weights []unitWeight
	switch {
	case isConquered && isCityNextToSea:
		weights = []unitWeight{
			{Tank, 1},
			{Fighter, 1},
			{Bomber, 1},
			{Transport, 1},
			{Destroyer, 2},
			{Submarine, 2},
			{Carrier, 2},
			{Battleship, 3},
		}
	case isConquered && !isCityNextToSea && tankCount >= 10:
		weights = []unitWeight{
			{Tank, 1},
			{Fighter, 1},
			{Bomber, 2},
		}
	case isConquered && !isCityNextToSea && tankCount < 10:
		weights = []unitWeight{
			{Tank, 5},
			{Fighter, 1},
			{Bomber, 1},
		}
	case !isConquered && isCityNextToSea && tankCount >= 10:
		weights = []unitWeight{
			{Tank, 1},
			{Fighter, 2},
			{Destroyer, 3},
		}
	case !isConquered && isCityNextToSea && tankCount < 10:
		weights = []unitWeight{
			{Tank, 3},
			{Fighter, 3},
			{Destroyer, 3},
		}
	default:
		weights = []unitWeight{
			{Tank, 7},
			{Fighter, 3},
		}
	}

	return getRandomUnit(weights)
}

// getUnitCount return a count of units of a given type within a islandMap for a player
func (g *GameBoard) getUnitCount(unitType UnitType, islandMap []Coordinate, player int) int {
	count := 0
	for _, coord := range islandMap {
		for _, unit := range g.Units {
			if unit.Type == unitType && unit.Player == player && unit.PositionX == coord.PositionX && unit.PositionY == coord.PositionY {
				count++
			}
		}
	}
	return count
}

// getRandomUnit calculates the total weight and selects a unit type based on these weights
func getRandomUnit(weights []unitWeight) UnitType {
	totalWeight := 0
	for _, w := range weights {
		totalWeight += w.weight
	}
	randomNum := rand.Intn(totalWeight) + 1
	currentWeight := 0
	for _, w := range weights {
		currentWeight += w.weight
		if randomNum <= currentWeight {
			return w.unit
		}
	}
	return Tank // Default to Tank if weights are not configured correctly
}

// hasPlayerWon checks if the specified player has won the game.
func (g *GameBoard) hasPlayerWon(playerID int) bool {
	// Check if any city is occupied by a different player
	for _, city := range g.Cities {
		if city.OccupyingPlayer == Unoccupied {
			return false // Player has not won if any city is unoccupied
		} else if city.OccupyingPlayer == OccupiedByPlayer1 && playerID != 1 {
			return false // Player has not won if any city is occupied by another player
		} else if city.OccupyingPlayer == OccupiedByPlayer2 && playerID != 2 {
			return false // Player has not won if any city is occupied by another player
		}
	}

	// Check if there are no enemy units on the board
	for _, unit := range g.Units {
		if unit.Player != playerID {
			return false // Player has not won if any enemy unit exists
		}
	}

	// If all cities are occupied by the player and there are no enemy units, the player has won
	return true
}

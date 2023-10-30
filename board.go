package main

import (
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
}

// Cell struct represents a cell on the game board.
type Cell struct {
	IsLand  bool // true for land, false for sea
	IsFog   bool // true for fog of war, false if visible
	HasCity bool // true if the cell has a city, false otherwise
}

// BoardCoordinate struct represents an X, Y, position on the game board.
type BoardCoordinate struct {
	PositionX int
	PositionY int
}

// NewGameBoard creates a new game board with the specified number of rows and columns.
func NewGameBoard(rows, columns int) *GameBoard {
	grid := make([][]Cell, rows)
	//for i := range grid {
	//	grid[i] = make([]Cell, columns)
	//}
	for i := range grid {
		row := make([]Cell, columns)
		for j := range row {
			row[j].IsFog = true // Initialize Fog to true for all cells
			// the default value of IsLand is 0 (false) by default, which represents Sea
		}
		grid[i] = row
	}

	return &GameBoard{
		Rows:    rows,
		Columns: columns,
		Grid:    grid,
	}
}

// GenerateRandomIslands generates random oval-shaped islands on the game board.
func (g *GameBoard) GenerateRandomIslands(numIslands int) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numIslands; i++ {
		centerRow := rand.Intn(g.Rows)
		centerCol := rand.Intn(g.Columns)
		radiusX := rand.Intn(6) + 2 // Random oval radius between 2 and 5 cells
		radiusY := rand.Intn(6) + 2 // Random oval radius between 2 and 5 cells

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
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numCities; i++ {
		for {
			row := rand.Intn(g.Rows)
			col := rand.Intn(g.Columns)

			// Check if the cell is land and does not have neighboring cities
			excludeTargetCell := false
			if g.Grid[row][col].IsLand && !g.HasNeighboringCity(row, col, excludeTargetCell) {
				//g.Grid[row][col].Land = true
				g.Grid[row][col].HasCity = true
				city := NewCity(row, col)
				city.IsCityNextToSea = g.IsCityNextToSea(city.PositionX, city.PositionY)
				g.Cities = append(g.Cities, *city)
				break
			}
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
				grid[i][j]="?"
			} else {
				if g.Grid[i][j].HasCity {
					grid[i][j]="C"
				} else if g.Grid[i][j].IsLand {
					grid[i][j]="L"
				} else {
					grid[i][j]="S"
				}
			}
		}
	}
	return grid
}

func (g *GameBoard) DoPlayerTurnAI(player int) {
	var activeUnit *Unit
	for {
		activeUnit = g.getActiveUnitForPlayer(player)
		if activeUnit == nil {
			break // No more active units for the player
		}
		// Process the active unit here
		g.runUnitAI(activeUnit)
	}
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

// getPossibleMoves returns possible moves for the given unit.
//
// priorities:
// should attack enemy unit if nearby
// should attack enemy city if nearby
// should attack unoccupied city if nearby
// should move to clear fog of war if nearby
func (g *GameBoard) getPossibleMoves(unit *Unit) []BoardCoordinate {
	var moves []BoardCoordinate
	var enemyUnits []BoardCoordinate
	var enemyCities []BoardCoordinate
	var unoccupiedCities []BoardCoordinate
	var fogOfWar []BoardCoordinate
	var randomMoves []BoardCoordinate

	// Check neighboring cells and add valid moves to the list
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue // Skip the current cell
			}
			newRow, newCol := unit.PositionX+i, unit.PositionY+j
			if newRow >= 0 && newRow < g.Rows && newCol >= 0 && newCol < g.Columns {
                defender := g.getUnitAtCoordinates(BoardCoordinate{newRow, newCol}, unit.Player)
	            if defender != nil {
	                enemyUnits = append(enemyUnits, BoardCoordinate{newRow, newCol})
	            }
				if g.Grid[newRow][newCol].HasCity {
				    city := g.getCityAtCoordinates(BoardCoordinate{newRow, newCol})
				    if city.OccupyingPlayer == Unoccupied {
				        unoccupiedCities = append(unoccupiedCities, BoardCoordinate{newRow, newCol})
				    } else if city.OccupyingPlayer == OccupiedByPlayer1 && unit.Player!=1 {
                        enemyCities = append(enemyCities, BoardCoordinate{newRow, newCol})
                    } else if city.OccupyingPlayer == OccupiedByPlayer2 && unit.Player!=2 {
                        enemyCities = append(enemyCities, BoardCoordinate{newRow, newCol})
                    }
                }
				if g.Grid[newRow][newCol].IsFog {
				    fogOfWar = append(fogOfWar, BoardCoordinate{newRow, newCol})
				}
				randomMoves = append(randomMoves, BoardCoordinate{newRow, newCol})
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
	    for _, randomMove := range randomMoves {
	        moves = append(moves, randomMove)
	    }
	}
	return moves
}

// attemptMoveTo attempts to move the unit to the destination coordinates
func (g *GameBoard) attemptMoveTo(destinationCoordinate BoardCoordinate, unit *Unit) {
	radius := 1
	g.clearFogOfWarAroundCoordinate(destinationCoordinate, radius)
	defender := g.getUnitAtCoordinates(destinationCoordinate, unit.Player)
	if defender != nil {
		g.resolveUnitAttack(unit, defender)
	} else if g.Grid[destinationCoordinate.PositionX][destinationCoordinate.PositionY].HasCity && unit.CanMoveOnLand {
		defender := g.getCityAtCoordinates(destinationCoordinate)
		g.resolveCityAttack(unit, defender)
	} else if g.Grid[destinationCoordinate.PositionX][destinationCoordinate.PositionY].IsLand && unit.CanMoveOnLand {
		unit.MoveTo(destinationCoordinate)
	} else if !g.Grid[destinationCoordinate.PositionX][destinationCoordinate.PositionY].IsLand && unit.CanMoveOnWater {
		unit.MoveTo(destinationCoordinate)
	} else {
		fmt.Println("attemptMoveTo() illegal move!")
	}
}

// clearFogOfWarAroundCoordinate clears the fog of war around the specified coordinates within a given radius.
func (g *GameBoard) clearFogOfWarAroundCoordinate(boardCoordinate BoardCoordinate, radius int) {
	for i := boardCoordinate.PositionX - radius; i <= boardCoordinate.PositionX+radius; i++ {
		for j := boardCoordinate.PositionY - radius; j <= boardCoordinate.PositionY+radius; j++ {
			if i >= 0 && i < g.Rows && j >= 0 && j < g.Columns {
				g.Grid[i][j].IsFog = false
			}
		}
	}
}

// getUnitAtCoordinates retrieves an enemy unit at the specified coordinates.
func (g *GameBoard) getUnitAtCoordinates(boardCoordinate BoardCoordinate, attackingPlayer int) *Unit {
	for i := range g.Units {
		if g.Units[i].PositionX == boardCoordinate.PositionX && g.Units[i].PositionY == boardCoordinate.PositionY && g.Units[i].Player != attackingPlayer {
			return &g.Units[i]
		}
	}
	return nil
}

// getCityAtCoordinates retrieves a city at the specified coordinates.
func (g *GameBoard) getCityAtCoordinates(boardCoordinate BoardCoordinate) *City {
	for i := range g.Cities {
		if g.Cities[i].PositionX == boardCoordinate.PositionX && g.Cities[i].PositionY == boardCoordinate.PositionY {
			return &g.Cities[i]
		}
	}
	return nil
}

// resolveCityAttack determines the outcome of an attack between an attacking unit and a defending city.
func (g *GameBoard) resolveCityAttack(attacker *Unit, defender *City) {
	attacker.MovesLeftThisDay--
	if attacker.CanFly {
		attacker.Fuel--
	}
	// Seed the random number generator with the current timestamp
	rand.Seed(time.Now().UnixNano())
	attackOutcome := false
	if rand.Intn(2) == 0 { // 50% probability.
		attackOutcome = true
	}
	if attackOutcome && attacker.Strength >= defender.Strength {
		// Apply damage to the defender's strength
		defender.Strength--
		// Check if the defender is destroyed
		if defender.Strength <= 0 {
			// Defender is conquered, change OccupyingPlayer
			if attacker.Player == 1 {
				defender.OccupyingPlayer = OccupiedByPlayer1
			} else {
				defender.OccupyingPlayer = OccupiedByPlayer2
			}
			defender.Strength = NewCityStrength
			defender.ManufacturingUnit = Blank // TODO: if player is computer, decide what to manufacture
			defender.DaysUntilUnitReady = 0
			// Attacker is destroyed when it conquers a city
			g.removeUnit(attacker)
		}
	} else {
		// Apply damage to the attacker's strength
		attacker.Strength--
		// Check if the attacker is destroyed
		if attacker.Strength <= 0 {
			// Attacker is destroyed, remove it from the game board
			g.removeUnit(attacker)
		}
	}
}

// resolveUnitAttack determines the outcome of an attack between an attacking unit and a defending unit.
func (g *GameBoard) resolveUnitAttack(attacker, defender *Unit) {
	attacker.MovesLeftThisDay--
	if attacker.CanFly {
		attacker.Fuel--
	}
	// Seed the random number generator with the current timestamp
	rand.Seed(time.Now().UnixNano())
	attackOutcome := false
	if rand.Intn(2) == 0 { // 50% probability.
		attackOutcome = true
	}
	if attackOutcome && attacker.Strength >= defender.Strength {
		// Apply damage to the defender's strength
		defender.Strength--
		// Check if the defender is destroyed
		if defender.Strength <= 0 {
			// Defender is destroyed, remove it from the game board
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

/*
// runTankAI implements the AI logic for the tank unit.
func (g *GameBoard) runTankAI() {
	rand.Seed(time.Now().UnixNano())

	// Sample AI logic
	// In this example, the tank moves randomly to neighboring cells
	if len(g.Units) > 0 {
		tankIndex := rand.Intn(len(g.Units))
		tank := g.Units[tankIndex]

		possibleMoves := g.getPossibleMoves(tank)
		if len(possibleMoves) > 0 {
			move := possibleMoves[rand.Intn(len(possibleMoves))]
			g.Units[tankIndex] = move // Update tank's position
		}
	}
}

// getPossibleMoves returns possible moves for the given unit.
func (g *GameBoard) getPossibleMoves(unit Unit) []Unit {
	var moves []Unit

	// Check neighboring cells and add valid moves to the list
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue // Skip the current cell
			}
			newRow, newCol := unit.Row+i, unit.Col+j
			if newRow >= 0 && newRow < g.Rows && newCol >= 0 && newCol < g.Columns {
				moves = append(moves, Unit{Row: newRow, Col: newCol})
			}
		}
	}

	return moves
}
*/

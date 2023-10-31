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

func (g *GameBoard) DoPlayerTurnAI(player int) {
	var activeUnit *Unit
	for {
		activeUnit = g.getActiveUnitForPlayer(player)
		if activeUnit == nil {
			break // No more active units for the player
		}
		// Process the active unit here
		g.runUnitAI(activeUnit)
		if g.hasPlayerWon(player) {
			break // the player has won
		}
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
		fmt.Println("found enemyUnits")
		moves = append(moves, enemyUnits[0])
	} else if len(enemyCities) > 0 {
		fmt.Println("found enemyCities")
		moves = append(moves, enemyCities[0])
	} else if len(unoccupiedCities) > 0 {
		fmt.Println("found unoccupiedCities")
		moves = append(moves, unoccupiedCities[0])
	} else if len(fogOfWar) > 0 {
		fmt.Println("found fogOfWar")
		moves = append(moves, fogOfWar[0])
	} else {
		fmt.Println("else")
		islandMap := g.getIslandMap(*unit)
		fmt.Printf("got islandMap %v\n", islandMap)
		isConquered := g.isIslandConquered(islandMap, unit.Player)
		fmt.Println("got isConquered")
		if isConquered {
			stagingPoint := g.getIsIslandCityNextToSea(islandMap)
			fmt.Printf("got stagingPoint %v\n", stagingPoint)
			if stagingPoint != nil {
				pathToStagingPoint := g.FindPath(*stagingPoint, unit)
				fmt.Printf("got pathToStagingPoint %v\n", pathToStagingPoint)
				if pathToStagingPoint != nil {
					firstStepOnPathTowardsStagingPoint := getSecondCoordinate(pathToStagingPoint)
					fmt.Printf("got firstStepOnPathTowardsStagingPoint %v\n", firstStepOnPathTowardsStagingPoint)
					if firstStepOnPathTowardsStagingPoint != nil {
						moves = append(moves, *firstStepOnPathTowardsStagingPoint)
					}
				}
			}
		}
	}
	if len(moves) == 0 {
		fmt.Println("use randomMoves")
		for _, randomMove := range randomMoves {
			moves = append(moves, randomMove)
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

// attemptMoveTo attempts to move the unit to the destination coordinates
func (g *GameBoard) attemptMoveTo(destinationCoordinate Coordinate, unit *Unit) {
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

// getAttackOutcome decides an attack outcome based on chance
func (g *GameBoard) getAttackOutcome() bool {
	// Seed the random number generator with the current timestamp
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 { // 50% probability.
		return true
	}
	return false
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
	attacker.MovesLeftThisDay--
	if attacker.CanFly {
		attacker.Fuel--
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
func (g *GameBoard) resolveUnitAttack(attacker, defender *Unit, attackOutcome bool) {
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

// getIslandMap returns a slice of coordinates representing the island connected to the given unit.
func (g *GameBoard) getIslandMap(unit Unit) []Coordinate {
	visited := make(map[Coordinate]bool)
	islandMap := make([]Coordinate, 0)

	// Define a recursive flood fill function to explore land cells
	var floodFill func(x, y int)
	floodFill = func(x, y int) {
		// Check if the cell is within the grid boundaries and is a land cell
		if x >= 0 && x < g.Rows && y >= 0 && y < g.Columns && g.Grid[y][x].IsLand {
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

	// Start the flood fill from the unit's position
	floodFill(unit.PositionX, unit.PositionY)

	return islandMap
}

// isIslandConquered checks if all cities on the island represented by coordinates are occupied by the same player.
func (g *GameBoard) isIslandConquered(islandMap []Coordinate, playerID int) bool {
	for _, coord := range islandMap {
		city := g.getCityAtCoordinates(coord)
		if city != nil {
			fmt.Printf("found city %d, %d\n", city.PositionX, city.PositionY)
			if city.OccupyingPlayer == Unoccupied {
				fmt.Println("...city is Unoccupied")
				return false // Island is not conquered if any city is unoccupied
			} else if city.OccupyingPlayer == OccupiedByPlayer2 && playerID == 1 {
				fmt.Println("...city is OccupiedByPlayer2")
				return false // Island is not conquered if any city is not occupied by the player
			} else if city.OccupyingPlayer == OccupiedByPlayer1 && playerID == 2 {
				fmt.Println("...city is OccupiedByPlayer1")
				return false // Island is not conquered if any city is not occupied by the player
			}
		}
	}
	return true // Island is conquered if all cities are occupied by the player
}

// getIsIslandCityNextToSea returns coordinate of city which is next to sea.
func (g *GameBoard) getIsIslandCityNextToSea(islandMap []Coordinate) *Coordinate {
	for _, coord := range islandMap {
		city := g.getCityAtCoordinates(coord)
		if city != nil && city.IsCityNextToSea {
			fmt.Printf("found island city next to sea %d, %d\n", city.PositionX, city.PositionY)
			return &Coordinate{PositionX: city.PositionX, PositionY: city.PositionY}
		}
	}
	fmt.Println("did not find island city next to sea, returning nil")
	return nil // city next to sea not found on island
}

// FindPath finds a path for the unit to reach the target coordinate on the grid.
func (g *GameBoard) FindPath(target Coordinate, unit *Unit) []Coordinate {
	fmt.Printf("in FindPath,\n")
	fmt.Printf("FindPath unit %d, %d\n", unit.PositionX, unit.PositionY)
	fmt.Printf("FindPath target %d, %d\n", target.PositionX, target.PositionY)
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

	fmt.Printf("No path found, returning nil\n")
	// No path found
	return nil
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

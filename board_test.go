package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestClearFogOfWarAroundCoordinate(t *testing.T) {

	want := [][]string{
		{"?", "?", "?", "?", "?", "?", "?", "?"},
		{"?", "?", "?", "L", "L", "L", "?", "?"},
		{"?", "?", "?", "L", "C", "L", "?", "?"},
		{"?", "?", "?", "L", "L", "L", "?", "?"},
	}

	rows, columns := 4, 8 // x, y (horizontal, vertical)
	board := NewGameBoard(rows, columns)
	// Initialize Land to true for all cells
	board.IterateGrid(func(row, col int, cell *Cell) {
		board.Grid[row][col].IsLand = true
	})
	// add a city
	cityRow := 2
	cityCol := 4
	board.Grid[cityRow][cityCol].HasCity = true
	showFogOfWar := false
	//board.Print(showFogOfWar)
	showFogOfWar = true
	//board.Print(showFogOfWar)
	Coordinate := Coordinate{2, 4}
	radius := 1
	board.clearFogOfWarAroundCoordinate(Coordinate, radius)
	//board.Print(showFogOfWar)
	got := board.printToSlice(showFogOfWar)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("name: ClearFogOfWarAroundCoordinate, got = %v; want %v", got, want)
	}
}

func TestHasNeighboringCity(t *testing.T) {
	rows, columns := 4, 8 // x, y (horizontal, vertical)
	board := NewGameBoard(rows, columns)
	// Initialize Land to true for all cells
	board.IterateGrid(func(row, col int, cell *Cell) {
		board.Grid[row][col].IsLand = true
	})
	cityRow := 2
	cityCol := 4
	board.Grid[cityRow][cityCol].HasCity = true
	//showFogOfWar := false
	//board.Print(showFogOfWar)
	want := false
	excludeTargetCell := true
	got := board.HasNeighboringCity(cityRow, cityCol, excludeTargetCell)
	if got != want {
		t.Errorf("name: no other city on map, got = %t; want %t", got, want)
	}

	board.Grid[cityRow][cityCol+1].HasCity = true
	//board.Print(showFogOfWar)
	want = true
	got = board.HasNeighboringCity(cityRow, cityCol, excludeTargetCell)
	if got != want {
		t.Errorf("name: 2nd city at col+1, got = %t; want %t", got, want)
	}

	board.Grid[cityRow][cityCol+1].HasCity = false
	board.Grid[cityRow-1][cityCol+1].HasCity = true
	//board.Print(showFogOfWar)
	want = true
	got = board.HasNeighboringCity(cityRow, cityCol, excludeTargetCell)
	if got != want {
		t.Errorf("name: 2nd city at row-1 col+1, got = %t; want %t", got, want)
	}

	board.Grid[cityRow][cityCol+1].HasCity = false
	board.Grid[cityRow-1][cityCol+1].HasCity = false
	board.Grid[cityRow-2][cityCol+1].HasCity = true
	//board.Print(showFogOfWar)
	want = false
	got = board.HasNeighboringCity(cityRow, cityCol, excludeTargetCell)
	if got != want {
		t.Errorf("name: 2nd city at row-2 col+1, got = %t; want %t", got, want)
	}
}

func TestIsCityNextToSea(t *testing.T) {
	rows, columns := 4, 8 // x, y (horizontal, vertical)
	board := NewGameBoard(rows, columns)
	// Initialize Land to true for all cells
	board.IterateGrid(func(row, col int, cell *Cell) {
		board.Grid[row][col].IsLand = true
	})
	cityRow := 2
	cityCol := 4
	board.Grid[cityRow][cityCol].HasCity = true
	//showFogOfWar := false
	//board.Print(showFogOfWar)
	want := false
	got := board.IsCityNextToSea(cityRow, cityCol)
	if got != want {
		t.Errorf("TestIsCityNextToSea() got = %t; want %t", got, want)
	}

	board.Grid[cityRow][cityCol+1].IsLand = false
	//board.Print(showFogOfWar)
	want = true
	got = board.IsCityNextToSea(cityRow, cityCol)
	if got != want {
		t.Errorf("TestIsCityNextToSea() got = %t; want %t", got, want)
	}

	board.Grid[cityRow][cityCol+1].IsLand = true
	board.Grid[cityRow-1][cityCol+1].IsLand = false
	//board.Print(showFogOfWar)
	want = true
	got = board.IsCityNextToSea(cityRow, cityCol)
	if got != want {
		t.Errorf("TestIsCityNextToSea() got = %t; want %t", got, want)
	}

	board.Grid[cityRow][cityCol+1].IsLand = true
	board.Grid[cityRow-1][cityCol+1].IsLand = true
	board.Grid[cityRow-2][cityCol+1].IsLand = false
	//board.Print(showFogOfWar)
	want = false
	got = board.IsCityNextToSea(cityRow, cityCol)
	if got != want {
		t.Errorf("TestIsCityNextToSea() got = %t; want %t", got, want)
	}
}

func TestRemoveUnit(t *testing.T) {
	// Create a GameBoard with some initial units
	initialUnits := []Unit{
		{PositionX: 1, PositionY: 1},
		{PositionX: 2, PositionY: 2},
		{PositionX: 3, PositionY: 3},
	}
	gameBoard := &GameBoard{Units: initialUnits}

	// Define the unit to be removed
	unitToRemove := &Unit{PositionX: 2, PositionY: 2}

	// Call the removeUnit function
	gameBoard.removeUnit(unitToRemove)

	// Define the expected units after removal
	expectedUnits := []Unit{
		{PositionX: 1, PositionY: 1},
		{PositionX: 3, PositionY: 3},
	}

	// Check if the game board's Units slice matches the expected units
	if !reflect.DeepEqual(gameBoard.Units, expectedUnits) {
		t.Errorf("Test failed: Expected units %+v, got %+v", expectedUnits, gameBoard.Units)
	}
}

// Test case 1: There are enemy units, move towards the first enemy unit
func TestGetPossibleMovesTestCase1(t *testing.T) {
	// Mock GameBoard with grid cells representing land and sea
	gameBoard := GameBoard{
		Rows:    3,
		Columns: 3,
		Grid: [][]Cell{
			{{IsLand: true}, {IsLand: true}, {IsLand: true}},
			{{IsLand: true}, {IsLand: true}, {IsLand: true}},
			{{IsLand: true}, {IsLand: true}, {IsLand: true}},
		},
		Cities: []City{},
	}

	showFogOfWar := false
	gameBoard.Print(showFogOfWar)

	unit := &Unit{
		PositionX:        2,
		PositionY:        1,
		CanMoveOnLand:    true,
		CanMoveOnWater:   false,
		CanFly:           false,
		Player:           1,
		MovesLeftThisDay: 2,
	}

	enemyUnit := &Unit{
		PositionX:        2,
		PositionY:        2,
		CanMoveOnLand:    true,
		CanMoveOnWater:   false,
		CanFly:           false,
		Player:           2,
		MovesLeftThisDay: 2,
	}

	gameBoard.Units = append(gameBoard.Units, *unit)
	gameBoard.Units = append(gameBoard.Units, *enemyUnit)

	possibleMoves := gameBoard.getPossibleMoves(unit)
	expectedMoves := []Coordinate{{2, 2}}
	if !slicesEqual(possibleMoves, expectedMoves) {
		t.Errorf("Expected moves: %v, but got: %v", expectedMoves, possibleMoves)
	}
}

// Test case 2: No enemy units, but enemy cities exist, move towards the first enemy city
func TestGetPossibleMovesTestCase2(t *testing.T) {
	// Mock GameBoard with grid cells representing land and sea
	gameBoard := GameBoard{
		Rows:    3,
		Columns: 3,
		Grid: [][]Cell{
			{{IsLand: true, HasCity: true}, {IsLand: true}, {IsLand: true}},
			{{IsLand: true}, {IsLand: true}, {IsLand: true}},
			{{IsLand: true}, {IsLand: true}, {IsLand: true}},
		},
		Cities: []City{
			{PositionX: 0, PositionY: 0, OccupyingPlayer: OccupiedByPlayer2},
		},
	}

	showFogOfWar := false
	gameBoard.Print(showFogOfWar)

	unit := &Unit{
		PositionX:        0,
		PositionY:        1,
		CanMoveOnLand:    true,
		CanMoveOnWater:   false,
		CanFly:           false,
		Player:           1,
		MovesLeftThisDay: 2,
	}

	gameBoard.Units = append(gameBoard.Units, *unit)

	possibleMoves := gameBoard.getPossibleMoves(unit)
	expectedMoves := []Coordinate{{0, 0}}
	if !slicesEqual(possibleMoves, expectedMoves) {
		t.Errorf("Expected moves: %v, but got: %v", expectedMoves, possibleMoves)
	}
}

// Test case 3: No enemy units or cities, move towards the first unoccupied city
func TestGetPossibleMovesTestCase3(t *testing.T) {
	// Mock GameBoard with grid cells representing land and sea
	gameBoard := GameBoard{
		Rows:    3,
		Columns: 3,
		Grid: [][]Cell{
			{{IsLand: true, HasCity: true}, {IsLand: true}, {IsLand: true}},
			{{IsLand: true}, {IsLand: true}, {IsLand: true}},
			{{IsLand: true}, {IsLand: true}, {IsLand: true}},
		},
		Cities: []City{
			{PositionX: 0, PositionY: 0, OccupyingPlayer: Unoccupied},
		},
	}

	showFogOfWar := false
	gameBoard.Print(showFogOfWar)

	unit := &Unit{
		PositionX:        0,
		PositionY:        1,
		CanMoveOnLand:    true,
		CanMoveOnWater:   false,
		CanFly:           false,
		Player:           1,
		MovesLeftThisDay: 2,
	}

	gameBoard.Units = append(gameBoard.Units, *unit)

	possibleMoves := gameBoard.getPossibleMoves(unit)
	expectedMoves := []Coordinate{{0, 0}}
	if !slicesEqual(possibleMoves, expectedMoves) {
		t.Errorf("Expected moves: %v, but got: %v", expectedMoves, possibleMoves)
	}
}

// Test case 4: No enemy units, cities, or unoccupied cities, move towards the fog of war cell
func TestGetPossibleMovesTestCase4(t *testing.T) {
	// Mock GameBoard with grid cells representing land and sea
	gameBoard := GameBoard{
		Rows:    3,
		Columns: 3,
		Grid: [][]Cell{
			{{IsLand: true, IsFog: true}, {IsLand: true}, {IsLand: true}},
			{{IsLand: true}, {IsLand: true}, {IsLand: true}},
			{{IsLand: true}, {IsLand: true}, {IsLand: true}},
		},
		Cities: []City{},
	}

	showFogOfWar := false
	gameBoard.Print(showFogOfWar)

	unit := &Unit{
		PositionX:        0,
		PositionY:        1,
		CanMoveOnLand:    true,
		CanMoveOnWater:   false,
		CanFly:           false,
		Player:           1,
		MovesLeftThisDay: 2,
	}

	gameBoard.Units = append(gameBoard.Units, *unit)

	possibleMoves := gameBoard.getPossibleMoves(unit)
	expectedMoves := []Coordinate{{0, 0}}
	if !slicesEqual(possibleMoves, expectedMoves) {
		t.Errorf("Expected moves: %v, but got: %v", expectedMoves, possibleMoves)
	}
}

// Test case 5: Island is conquered, move towards the staging point
func TestGetPossibleMovesTestCase5(t *testing.T) {
	// Mock GameBoard with grid cells representing land and sea
	gameBoard := GameBoard{
		Rows:    3,
		Columns: 3,
		Grid: [][]Cell{
			{{IsLand: true}, {IsLand: false}, {IsLand: false}},
			{{IsLand: true}, {IsLand: false}, {IsLand: false}},
			{{IsLand: true}, {IsLand: true}, {IsLand: true, HasCity: true}},
		},
		Cities: []City{
			{PositionX: 2, PositionY: 2, OccupyingPlayer: OccupiedByPlayer1, IsCityNextToSea: true},
		},
	}

	showFogOfWar := false
	gameBoard.Print(showFogOfWar)

	unit := &Unit{
		PositionX:        0,
		PositionY:        0,
		CanMoveOnLand:    true,
		CanMoveOnWater:   false,
		CanFly:           false,
		Player:           1,
		MovesLeftThisDay: 2,
	}

	fmt.Printf("unit %d, %d\n", unit.PositionX, unit.PositionY)

	gameBoard.Units = append(gameBoard.Units, *unit)

	possibleMoves := gameBoard.getPossibleMoves(unit)
	expectedMoves := []Coordinate{{1, 0}}
	if !slicesEqual(possibleMoves, expectedMoves) {
		t.Errorf("Expected moves: %v, but got: %v", expectedMoves, possibleMoves)
	}
}

func slicesEqual(slice1, slice2 []Coordinate) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func TestGetIslandMap(t *testing.T) {
	// Mock GameBoard with grid cells representing land and sea
	board := GameBoard{
		Rows:    3,
		Columns: 3,
		Grid: [][]Cell{
			{{IsLand: false}, {IsLand: true}, {IsLand: true}},
			{{IsLand: true}, {IsLand: true}, {IsLand: true}},
			{{IsLand: true}, {IsLand: true}, {IsLand: false}},
		},
	}

	//showFogOfWar := false
	//board.Print(showFogOfWar)

	// Mock unit's position
	unit := Unit{PositionX: 2, PositionY: 1}

	// Expected island map for the given unit's position
	expectedIslandMap := []Coordinate{
		{PositionX: 2, PositionY: 1},
		{PositionX: 1, PositionY: 0},
		{PositionX: 0, PositionY: 1},
		{PositionX: 0, PositionY: 2},
		{PositionX: 1, PositionY: 1},
		{PositionX: 1, PositionY: 2},
		{PositionX: 2, PositionY: 0},
	}

	// Call the function and check if the returned island map matches the expected one
	result := board.getIslandMap(Coordinate{unit.PositionX, unit.PositionY})
	if !reflect.DeepEqual(result, expectedIslandMap) {
		t.Errorf("Test failed: Expected island map %+v, got %+v", expectedIslandMap, result)
	}
}

func TestGetIsIslandCityNextToSea(t *testing.T) {
	// Mock GameBoard and cities
	gameBoard := &GameBoard{
		Rows:    5,
		Columns: 5,
		Grid:    nil, // set grid configuration here if needed
		Cities: []City{
			{PositionX: 0, PositionY: 0, IsCityNextToSea: false},
			{PositionX: 1, PositionY: 1, IsCityNextToSea: true},
			{PositionX: 2, PositionY: 2, IsCityNextToSea: true},
			{PositionX: 3, PositionY: 3, IsCityNextToSea: false},
			{PositionX: 4, PositionY: 4, IsCityNextToSea: true},
		},
	}

	// Test case 1: There is a city next to the sea, expect non-nil Coordinate
	islandMap := []Coordinate{{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}}
	result := gameBoard.getIsIslandCityNextToSea(islandMap)
	expectedResult := &Coordinate{PositionX: 1, PositionY: 1}
	assertCoordinatesEqual(t, result, expectedResult, "Test case 1")

	// Test case 2: There is no city next to the sea, expect nil Coordinate
	islandMap = []Coordinate{{0, 0}, {3, 3}}
	result = gameBoard.getIsIslandCityNextToSea(islandMap)
	assertCoordinatesEqual(t, result, nil, "Test case 2")
}

func assertCoordinatesEqual(t *testing.T, got, expected *Coordinate, testName string) {
	t.Helper()
	if (got == nil && expected != nil) || (got != nil && expected == nil) {
		t.Errorf("%s: Expected: %v, but got: %v", testName, expected, got)
		return
	}
	if got != nil && expected != nil && *got != *expected {
		t.Errorf("%s: Expected: %v, but got: %v", testName, expected, got)
	}
}

func TestIsIslandConquered(t *testing.T) {
	// Mock game board with cities
	gameBoard := GameBoard{
		Cities: []City{
			{PositionX: 1, PositionY: 1, OccupyingPlayer: OccupiedByPlayer1},
			{PositionX: 2, PositionY: 2, OccupyingPlayer: OccupiedByPlayer1},
			{PositionX: 3, PositionY: 3, OccupyingPlayer: OccupiedByPlayer1},
			{PositionX: 4, PositionY: 4, OccupyingPlayer: Unoccupied},
		},
	}

	// Test case 1: All cities are occupied by player 1, island is conquered
	islandMap := []Coordinate{
		{PositionX: 1, PositionY: 1},
		{PositionX: 2, PositionY: 2},
		{PositionX: 3, PositionY: 3},
	}
	if !gameBoard.isIslandConquered(islandMap, 1) {
		t.Errorf("Test case 1 failed: Island should be conquered by player 1")
	}

	// Test case 2: One city is unoccupied, island is not conquered
	islandMap = []Coordinate{
		{PositionX: 1, PositionY: 1},
		{PositionX: 2, PositionY: 2},
		{PositionX: 3, PositionY: 3},
		{PositionX: 4, PositionY: 4}, // Unoccupied city
	}
	if gameBoard.isIslandConquered(islandMap, 1) {
		t.Errorf("Test case 2 failed: Island should not be conquered due to unoccupied city")
	}

	// Test case 3: Player 1 tries to conquer an island occupied by player 2, should fail
	gameBoard.Cities[2].OccupyingPlayer = OccupiedByPlayer2 // Change a city to be occupied by player 2
	islandMap = []Coordinate{
		{PositionX: 1, PositionY: 1},
		{PositionX: 2, PositionY: 2},
		{PositionX: 3, PositionY: 3},
	}
	if gameBoard.isIslandConquered(islandMap, 1) {
		t.Errorf("Test case 3 failed: Island should not be conquered due to presence of enemy city")
	}
}

// TestHasPlayerWon checks the hasPlayerWon function.
func TestHasPlayerWon(t *testing.T) {
	// Mock GameBoard with cities and units for testing
	gameBoard := GameBoard{
		Cities: []City{
			{OccupyingPlayer: 1},
			{OccupyingPlayer: 1},
			{OccupyingPlayer: 1},
		},
		Units: []Unit{
			{Player: 1},
			{Player: 1},
		},
	}

	// Test case 1: Player 1 has won
	if !gameBoard.hasPlayerWon(1) {
		t.Errorf("Test case 1 failed: Player 1 should have won")
	}

	// Test case 2: Player 2 has not won as there is a city controlled by a different player
	if gameBoard.hasPlayerWon(2) {
		t.Errorf("Test case 2 failed: Player 2 should not have won")
	}

	// Test case 3: Player 3 has not won as there is an enemy unit
	if gameBoard.hasPlayerWon(3) {
		t.Errorf("Test case 3 failed: Player 3 should not have won")
	}
}

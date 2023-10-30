package main

import (
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
	cityRow := 2
	cityCol := 4

	board.Grid[cityRow][cityCol].HasCity = true
	showFogOfWar := false
	//board.Print(showFogOfWar)
	showFogOfWar = true
	//board.Print(showFogOfWar)
	boardCoordinate := BoardCoordinate{2, 4}
	radius := 1
	board.clearFogOfWarAroundCoordinate(boardCoordinate, radius)
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

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {

	rows, columns := 10, 20 // x, y (horizontal, vertical) (rows, columns)
	board := NewGameBoard(rows, columns)
	numIslands := 4
	board.GenerateRandomIslands(numIslands)
	numCities := 12
	board.AddCities(numCities)
	board.Player1 = NewPlayer("player 1", true)
	board.Player2 = NewPlayer("player 2", true)
	board.DayZero()
	for {
		if board.Day == 20 {
			fmt.Println("demo game cut short")
			break
		}
		board.NextDay()
		board.DoPlayerTurnAI(1)
		if board.hasPlayerWon(1) {
			break
		}
		board.DoPlayerTurnAI(2)
		if board.hasPlayerWon(2) {
			break
		}
	}
	fmt.Println("GAME OVER")
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

package main

import (
    "fmt"
    "os"
    "os/exec"
    "runtime"
)

// ANSI color codes
const (
    Reset  = "\033[0m"
    Red    = "\033[31m"
    Green  = "\033[32m"
    Yellow = "\033[33m"
    Blue   = "\033[34m"
    Purple = "\033[35m"
    Cyan   = "\033[36m"
    White  = "\033[37m"
)

func main() {

	rows, columns := 10, 20 // x, y (horizontal, vertical) (rows, columns)
	board := NewGameBoard(rows, columns)
	numIslands := 10
	board.GenerateRandomIslands(numIslands)
	numCities := 30
	board.AddCities(numCities)
	board.Player1 = NewPlayer("player 1", true)
	board.Player2 = NewPlayer("player 2", true)
    board.DayZero()
    for {
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

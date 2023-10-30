package main

// "fmt"

func main() {
	rows, columns := 25, 55 // x, y (horizontal, vertical) (rows, columns)
	board := NewGameBoard(rows, columns)

	numIslands := 10
	board.GenerateRandomIslands(numIslands)

	numCities := 10
	board.AddCities(numCities)

	//player1 := NewPlayer("player1", true)
	//player2 := NewPlayer("player2", true)

	//day := 1
	//board.DoPlayerTurnAI(1)
	//board.DoPlayerTurnAI(2)

	//day++
	//board.DoPlayerTurnAI(1)
	//board.DoPlayerTurnAI(2)

	//showFogOfWar := false
	//board.Print(showFogOfWar)

	/*
				// Example usage
				city := NewCity(0, 0)
				city.OccupyCity(1) // Occupied by player 1
				SetManufacturingUnit(Tank)

				// Simulate days
				for day := 1; day <= 5; day++ { // Simulating 5 days
					fmt.Printf("Day %d\n", day)

					// loop list of cities occupied by player 2

		                if city.ManufactureUnit() {
		                    // TODO: add unit to list of units for current player
		                    fmt.Printf("City manufactured a %s!\n", city.ManufacturingUnit)
		                } else {
		                    fmt.Println("City is still manufacturing a unit...")
		                }

		                // TODO: give player 2 control until they indicate "end turn"

					// loop list of cities occupied by player 1

		                if city.ManufactureUnit() {
		                    // TODO: add unit to list of units for current player
		                    fmt.Printf("City manufactured a %s!\n", city.ManufacturingUnit)
		                } else {
		                    fmt.Println("City is still manufacturing a unit...")
		                }

		                // TODO: give player 1 control until they indicate "end turn"

					//time.Sleep(1 * time.Second) // Simulate some time passing between turns
				}
	*/
}

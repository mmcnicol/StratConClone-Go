package main

// Player struct represents a player in the game
type Player struct {
	Name string
	IsAI bool
}

func NewPlayer(name string, isAI bool) *Player {
	return &Player{
		Name: name,
		IsAI: isAI,
	}
}

package server

// Piece of the map
type Zone struct {
}

// Hub
type GameHub struct {
	world   World
	Updates chan Ship
}

// NewGameHub returns a new game and world
func NewGameHub() GameHub {
	return GameHub{
		world:   RandomWorld(100),
		Updates: make(chan Ship),
	}
}

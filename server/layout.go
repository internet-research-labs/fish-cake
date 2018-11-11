package server

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Pos3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type CityLayout struct {
	Buildings []Position `json:"buildings"`
}

// The world is a torus
type World struct {
	Buildings []Position `json:"buildings"`
	Thickness float64    `json:"thickness"`
	Radius    float64    `json:"radius"`
}

// Coordinates
func (s *World) Coord(t, f float64) (float64, float64, float64) {
	return 0.0, 0.0, 0.0
}

package server

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
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

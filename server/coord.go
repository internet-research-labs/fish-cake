package server

// Cartesian-3 Coord
type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// Cartesian-3 Coord
type Pos3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

// Torus Surface Coord
type WorldCoord struct {
	Theta float64 `json:"theta"`
	Fi    float64 `json:"fi"`
}

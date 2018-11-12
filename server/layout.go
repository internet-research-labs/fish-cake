package server

import (
	"math"
	"math/rand"
	"time"
)

// How things are laid out
type CityLayout struct {
	Buildings []WorldCoord `json:"buildings"`
}

// Ship
type Ship struct {
	Id    uint64     `json:"id"`
	Coord WorldCoord `json:"coord"`
}

// The world is a torus
type World struct {
	Buildings []WorldCoord    `json:"buildings"`
	Ships     map[uint64]Ship `json:"ships"`
	Radius    float64         `json:"radius"`
	Thickness float64         `json:"thickness"`
}

func RandomBuildings(n int) []WorldCoord {

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	_ = r

	buildings := make([]WorldCoord, 0)

	/*
		for i := 0; i < n; i++ {

			p := WorldCoord{
				Theta: r.Float64() * 2 * math.Pi,
				Fi:    r.Float64() * 2 * math.Pi,
			}

			buildings[i] = p
		}
	*/

	for i := 0; i < 20; i++ {
		for j := 0; j < 10; j++ {
			p := WorldCoord{
				Theta: float64(i) / 20.0 * 2 * math.Pi,
				Fi:    float64(j) / 10.0 * 2 * math.Pi,
			}
			buildings = append(buildings, p)
		}
	}

	return buildings
}

// MakeRandomWorld returns a world structure
func RandomWorld(n int) World {

	ships := make(map[uint64]Ship)

	// World
	return World{
		Buildings: RandomBuildings(400),
		Ships:     ships,
		Radius:    15.0,
		Thickness: 1.0,
	}
}

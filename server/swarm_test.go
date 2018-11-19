package server

import (
	"log"
	"testing"
)

// TestSwarm
func TestSwarm(t *testing.T) {
	zone := NewSwiftZone()

	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			x := float64(i) * 0.1
			y := float64(j) * 0.1
			zone.Add(
				Vector3{x, y, 0.0},
				Vector3{0.0, 0.0, 1.0},
			)
		}
	}

	zone.Tick()

	// panic("Let's make this start iterating")
}
func TestGetNear(t *testing.T) {
	zone := NewSwiftZone()
	d := Vector3{0.0, 0.0, 1.0}

	zone.Add(
		Vector3{1.0, 0.0, 0.0},
		d,
	)
	zone.Add(
		Vector3{0.0, 0.0, 0.0},
		d,
	)
	zone.Add(
		Vector3{0.0, 0.0, -1.0},
		d,
	)

	results := zone.GetNear(Vector3{0.0, 0.0, 0.0}, 0.1)

	if len(results) != 1 {
		t.Error("Should return exactly one swift")
	}
}

func o_o() {
	log.Println("o_o")
}

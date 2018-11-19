package server

import (
	"testing"
)

// TestSwarm
func TestSwarm(t *testing.T) {
	zone := NewSwiftZone()

	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			x := float64(i) * 0.1
			z := float64(j) * 0.1
			zone.Add(Vector3{x, 0.0, z})
		}
	}

	zone.Tick()

	panic("Let's make this start iterating")
}

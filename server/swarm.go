package server

import (
	"log"
)

// Swift
type Swift struct {
	Id  uint
	Pos Vector3
	Dir Vector3
}

// SwiftMap is the standard collection for swifts
type SwiftMap map[uint]*Swift

// Zone constants for swifts
const (
	SPEED   = 0.05
	ATTRACT = 0.1
	REPULSE = 0.1
	ORIENT  = 0.1
)

// Zone
type SwiftZone struct {
	id      uint            `json:"id"`
	swifts  map[uint]*Swift `json:"swifts"`
	channel chan Swift      `json:"channel"`
}

// NewSwiftZone returns a new zone
func NewSwiftZone() SwiftZone {
	return SwiftZone{
		swifts:  make(map[uint]*Swift),
		channel: make(chan Swift),
	}
}

// Add swift to map
func (self *SwiftZone) Add(p, d Vector3) {
	swifty := Swift{
		Id:  self.id,
		Pos: p,
		Dir: Vector3{1.0, 0.0, 0.0},
	}
	self.swifts[self.id] = &swifty
	self.id += 1
}

// Update ...
func (self *SwiftZone) Tick() {
	for _, v := range self.swifts {
		v.Pos.X += 0.1
	}
}

// GetNear returns a map of swifts near a point with radius d
func (self *SwiftZone) GetNear(pos Vector3, d float64) SwiftMap {
	neighbors := make(SwiftMap)
	for _, v := range self.swifts {
		if Distance(v.Pos, pos) < d {
			neighbors[v.Id] = v
		}
	}
	return neighbors
}

func x_x() {
	log.Println("x_x")
}

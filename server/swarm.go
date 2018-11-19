package server

import (
	"log"
)

const (
	SPEED   = 0.05
	ATTRACT = 0.1
	REPULSE = 0.1
	ORIENT  = 0.1
)

type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Add returns a new vector3
func Add(a, b Vector3) Vector3 {
	return Vector3{
		a.X + b.X,
		a.Y + b.Y,
		a.Z + b.Z,
	}
}

// Swift
type Swift struct {
	Id  uint
	Pos Vector3
	Dir Vector3
}

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
func (self *SwiftZone) Add(p Vector3) {
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
		log.Println("x_x", v)
	}
}

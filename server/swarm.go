package server

import (
	"log"
	"time"
)

// Swift
type Swift struct {
	Id  uint    `json:"id"`
	Pos Vector3 `json:"pos"`
	Dir Vector3 `json:"dir"`
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

// SwiftZone keeps track of all the swifts, and exposes a channel of updates
type SwiftZone struct {
	id      uint            `json:"id"`
	swifts  map[uint]*Swift `json:"swifts"`
	channel chan SwiftMap   `json:"channel"`
	ticker  *time.Ticker
}

// NewSwiftZone returns a new zone
func NewSwiftZone() SwiftZone {
	return SwiftZone{
		swifts:  make(map[uint]*Swift),
		channel: make(chan SwiftMap),
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

// Start ticking
func (self *SwiftZone) Start() {
	log.Println("SwiftZone.Start()")
	self.ticker = time.NewTicker(13 * time.Millisecond)
	go func() {
		for _ = range self.ticker.C {
			self.tick()
		}
	}()
}

func (self *SwiftZone) Stop() {
	log.Println("SwiftZone.Stop()")
	self.ticker.Stop()
	self.ticker = nil
}

// Update ...
func (self *SwiftZone) tick() {
	for _, v := range self.swifts {
		v.Pos.Z += 0.05

		// Wrap boundaries
		switch {
		case v.Pos.X > 8.0:
			v.Pos.X = -8.0
		case v.Pos.Y > 8.0:
			v.Pos.Y = -8.0
		case v.Pos.Z > 8.0:
			v.Pos.Z = -8.0
		}
	}

	self.channel <- self.GetNear(Vector3{0., 0., 0.}, 16.0)
}

// GetNear returns a map of swifts near a point with radius d
func (self *SwiftZone) GetNear(pos Vector3, d float64) SwiftMap {
	neighbors := make(SwiftMap)
	for _, v := range self.swifts {
		if Distance(v.Pos, pos) <= d {
			neighbors[v.Id] = v
		}
	}
	return neighbors
}

func x_x() {
	log.Println("x_x")
}

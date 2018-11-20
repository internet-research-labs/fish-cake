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
		Dir: d,
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

// updateposition updates a position for a swift
// - we are essentially creating a force field
func (self *SwiftZone) updatePosition(swift *Swift) {
	neighbors := self.GetNear(swift.Pos, 0.0)
	pos := Vector3{0.0, 0.0, 0.0}
	for _, n := range neighbors {
		pos = Add(pos, n.Pos)
	}

	attractor := Scale(pos, 1.0/float64(len(neighbors)))
	direction := swift.Dir
	norm := NormL2(direction)

	_ = attractor

	swift.Pos = Add(
		swift.Pos,
		Scale(direction, 0.1/norm),
	)

	panic("This isn't working as expected")
}

// updateposition updates a position for a swift
func (self *SwiftZone) wrap(swift *Swift) {
	// Wrap boundaries
	switch {
	case swift.Pos.X > 8.0:
		swift.Pos.X = -8.0
	case swift.Pos.Y > 8.0:
		swift.Pos.Y = -8.0
	case swift.Pos.Z > 8.0:
		swift.Pos.Z = -8.0

	case swift.Pos.X < -8.0:
		swift.Pos.X = 8.0
	case swift.Pos.Y < -8.0:
		swift.Pos.Y = 8.0
	case swift.Pos.Z < -8.0:
		swift.Pos.Z = 8.0
	}
}

// tick updates every swift
func (self *SwiftZone) tick() {
	for _, s := range self.swifts {
		self.updatePosition(s)
		self.wrap(s)
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

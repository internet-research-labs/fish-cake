package server

import (
	"log"
	"math"
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
// XXX: These aren't used yet
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
func (self *SwiftZone) Start(n time.Duration) {
	log.Println("SwiftZone.Start()")
	self.ticker = time.NewTicker(n)
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

// GetForce returns ...
func GetForce(swift, actor *Swift) Vector3 {
	mag := Distance(swift.Pos, actor.Pos)
	// Determine how much this {swift} wants to move towards this {actor}
	dir := Sub(actor.Pos, swift.Pos)

	shove := 0.0001

	if mag > 3.0 {
		return Scale(dir, shove)
	}
	return Scale(dir, -shove)
}

// getMapOfNewSwifts returns a map of swifts
// NOTE: ...
func (self *SwiftZone) getUpdatedPositions() map[uint]Swift {
	swifts := make(map[uint]Swift)

	for id, swift := range self.swifts {
		influence := Vector3{0.0, 0.0, 0.0}
		neighbors := self.GetNear(swift.Pos, 0.707)

		// For everyone near enough to influence...
		// Let's compute the overall force
		for _, n := range neighbors {
			force := GetForce(swift, n)
			influence = Add(influence, force)
		}

		dir := Add(swift.Dir, influence)

		// Attach swift to it
		swifts[id] = Swift{
			Id:  id,
			Pos: Add(swift.Pos, dir),
			Dir: swift.Dir,
		}
	}
	return swifts
}

// updateposition updates a position for a swift
func (self *SwiftZone) wrap(swift *Swift) {
	SIZE := 8.0
	AROU := 2 * SIZE
	swift.Pos.X = math.Mod(swift.Pos.X+SIZE, AROU) - SIZE
	swift.Pos.Y = math.Mod(swift.Pos.Y+SIZE, AROU) - SIZE
	swift.Pos.Z = math.Mod(swift.Pos.Z+SIZE, AROU) - SIZE
}

// tick updates every swift
func (self *SwiftZone) tick() {
	updateMap := self.getUpdatedPositions()

	for id, s := range self.swifts {
		newSwift, found := updateMap[id]
		if found {
			s.Pos.X = newSwift.Pos.X
			s.Pos.Y = newSwift.Pos.Y
			s.Pos.Z = newSwift.Pos.Z
		}
		// self.updatePosition(s)
		self.wrap(s)
	}

	self.channel <- self.GetNear(Vector3{0., 0., 0.}, 16.0)
}

// GetNear returns a map of swifts near a point with radius d
func (self *SwiftZone) GetNear(pos Vector3, d float64) SwiftMap {
	return self.swifts
	/*
		neighbors := make(SwiftMap)
		for _, v := range self.swifts {
			if Distance(v.Pos, pos) <= d {
				neighbors[v.Id] = v
			}
		}
		return neighbors
	*/
}

func x_x() {
	log.Println("x_x")
}

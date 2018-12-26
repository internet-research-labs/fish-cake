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

func (self *Swift) Steal(rhs *Swift) {
	// Update position
	self.Pos.X = rhs.Pos.X
	self.Pos.Y = rhs.Pos.Y
	self.Pos.Z = rhs.Pos.Z
}

// SwiftMap is the standard collection for swifts
type SwiftMap map[uint]*Swift

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
	close(self.channel)
	self.ticker.Stop()
	self.ticker = nil
}

// GetForce returns ...
func GetForce(swift, actor *Swift) Vector3 {

	const (
		ATTRACT = 0.001
		REPULSE = 0.0014
	)
	d := Distance(swift.Pos, actor.Pos)
	dir := Sub(actor.Pos, swift.Pos)
	mag := ATTRACT/d/d - REPULSE/d/d/d

	if d < 0.01 {
		return Vector3{0.0, 0.0, 0.0}
	}

	return Scale(dir, mag)
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
		for n_id, n := range neighbors {
			if id == n_id {
				continue
			}
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

	// Get set of updated positions per id
	updateMap := self.getUpdatedPositions()

	//
	for id, s := range self.swifts {
		newSwift, found := updateMap[id]
		if found {
			s.Steal(&newSwift)
		}
		self.wrap(s)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered")
		}
	}()

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

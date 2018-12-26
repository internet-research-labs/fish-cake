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

	// swarming properties
	attraction float64
	repulsion  float64
	alignment  float64
}

// NewSwiftZone returns a new zone
func NewSwiftZone(a, b, c float64) SwiftZone {
	return SwiftZone{
		swifts:  make(map[uint]*Swift),
		channel: make(chan SwiftMap),

		// Swarm properties
		attraction: a,
		repulsion:  b,
		alignment:  c,
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

func (self *Vector3) Copy() Vector3 {
	return Vector3{self.X, self.Y, self.Z}
}

// getMapOfNewSwifts returns a map of swifts
// NOTE: ...
func (self *SwiftZone) getUpdatedPositions() map[uint]Swift {

	swifts := make(map[uint]Swift)

	for id, swift := range self.swifts {
		influence := Vector3{0.0, 0.0, 0.0}
		neighbors := self.GetNear(swift.Pos, 0.707)
		aligner := Vector3{}

		// For everyone near enough to influence...
		// Let's compute the overall force
		for n_id, n := range neighbors {
			if id == n_id {
				continue
			}

			d := Distance(swift.Pos, n.Pos)

			// Get overall desired position to get near
			force := GetForce(swift, n)
			influence = Add(influence, force)

			if d < 0.9 && d > 0.1 {
				aligner = Add(aligner, n.Dir)
			}
		}

		aligner = aligner
		dir := Add(aligner, influence)

		// dir := Add(influence, swift.Dir)
		// dir = Scale(dir, 0.5)

		// Attach swift to it
		swifts[id] = Swift{
			Id:  id,
			Pos: Add(swift.Pos, dir),
			Dir: dir,
		}
	}
	return swifts
}

// Wrapf is like a mod that shifts negative numbers to positive ones
func Wrapf(x, y float64) float64 {
	if y < 0.0 {
		return 0.0
	}
	for x < 0.0 {
		x += y
	}

	if x > y {
		return math.Mod(x, y)
	}

	return x
}

// updateposition updates a position for a swift
func (self *SwiftZone) wrap(swift *Swift) {
	LOW := -8.0
	HIGH := 8.0
	AROU := HIGH - LOW
	swift.Pos.X = Wrapf(swift.Pos.X-LOW, AROU) + LOW
	swift.Pos.Y = Wrapf(swift.Pos.Y-LOW, AROU) + LOW
	swift.Pos.Z = Wrapf(swift.Pos.Z-LOW, AROU) + LOW
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

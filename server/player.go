package server

import (
	"github.com/gorilla/websocket"
)

const (
	BLACK = 0
	RED   = 1
)

// Player has the info to contact a player and update its stuff
type Player struct {
	ip         string
	ship       Ship
	connection *websocket.Conn
	channel    chan Ship
	side       uint
}

func (self *Player) Tick() {
	self.ship.Coord.Fi += 0.1
	self.channel <- self.ship
}

// SendMessage sends a message over that players connection
func (self *Player) SendMessage(m interface{}) {
	self.connection.WriteMessage(1, EncodeWireMessage("bleep", m))
}

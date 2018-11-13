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

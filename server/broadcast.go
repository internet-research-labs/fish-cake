package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// NOTE: Not family connections
type ConnectionFamily struct {
	conns map[*websocket.Conn]struct{}
	mux   sync.Mutex
}

func (family *ConnectionFamily) Add(conn *websocket.Conn) {
	family.mux.Lock()
	defer family.mux.Unlock()

	family.conns[conn] = struct{}{}
}

func (family *ConnectionFamily) Remove(conn *websocket.Conn) {
	family.mux.Lock()
	defer family.mux.Unlock()

	delete(family.conns, conn)
}

func (family *ConnectionFamily) Write(msg []byte) {
	log.Printf("Writing out to %v connections", len(family.conns))
	for conn := range family.conns {
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

var genericUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 8,
	WriteBufferSize: 1024 * 8,
}

type Broadcaster struct {
	family   ConnectionFamily
	messages chan []byte
}

// Handler returns a websocket http handler  which receives any message sent to In().
func (caster *Broadcaster) Handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("Incoming connection!")

		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		caster.family.Add(conn)
		defer caster.family.Remove(conn)

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}
}

// Start begins processing inbound messages.
func (caster *Broadcaster) Start() {
	go func() {
		for message := range caster.messages {
			caster.family.Write(message)
		}
	}()
}

// In returns a write-only channel of messages to broadcast.
func (caster *Broadcaster) In() chan<- []byte {
	return caster.messages
}

// NewBroadcaster returns a reference to a Broadcaster with some nice defaults.
func NewBroadcaster() *Broadcaster {
	fam := ConnectionFamily{
		conns: make(map[*websocket.Conn]struct{}),
	}
	return &Broadcaster{fam, make(chan []byte, 100)}
}

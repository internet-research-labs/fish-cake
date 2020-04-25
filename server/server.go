package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	STATIC_DIR = "/static/"
	PORT       = "8080"
)

// Server holds local zone game information
type Server struct {
	templates  string
	zones      []Zone
	players    map[uint64]*Player
	bots       []uint
	world      World
	simulation SwiftZone
}

// SocketHandler returns a handler function for gorilla that adds a new ship
func (self *Server) SocketHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Upgrade to a websocket connection
		conn, _ := upgrader.Upgrade(w, r, nil)

		// Create a new random ship
		player := Player{
			ip:         "whatever",
			ship:       NewRandomShip(),
			connection: conn,
			channel:    make(chan Ship),
		}

		// Add player
		self.players[player.ship.Id] = &player

		// It doesn't have to be this hard
		defer func() {
			log.Println("Closing connection:", player.ship.Id)
			delete(self.players, player.ship.Id)
			conn.Close()
		}()

		log.Println("New connection:", player)

		// READ LOOP
		go func() {
			for {
				t, m, e := conn.ReadMessage()
				if e != nil {
					log.Println("Error:", e, t, m)
					close(player.channel)
					return
				} else {
					log.Println("Message:", t, m)
				}
			}
		}()

		info, _ := json.Marshal(struct {
			Type string
			Id   string
		}{
			Type: "you-are",
			Id:   fmt.Sprintf("%d", player.ship.Id),
		})

		conn.WriteMessage(1, []byte(info))

		//
		world_message, _ := json.Marshal(struct {
			Type  string
			World World
		}{
			Type:  "world",
			World: self.world,
		})

		// Send message to the world
		conn.WriteMessage(1, []byte(world_message))

		// Message
		for s := range player.channel {
			m, e := json.Marshal(s)
			if e != nil {
				conn.WriteMessage(1, []byte(m))
			} else {
				continue
			}
		}
	}
}

// NewRandomServer returns a new server with a random world
func NewRandomServer(n int, static string) Server {
	return Server{
		templates:  static,
		zones:      nil,
		players:    make(map[uint64]*Player),
		bots:       nil,
		world:      RandomWorld(n),
		simulation: *NewSwiftZone(0.001, 0.0014, 0.5),
	}
}

// upgrader adds websocket support
var upgrader = websocket.Upgrader{}

// Listen for connections
func (self *Server) Listen(port int) {
	ticker := time.NewTicker(23 * time.Millisecond)

	// StaticHandler responds to requests by path
	StaticHandler := func(res http.ResponseWriter, req *http.Request) {
		fname := path.Base(req.URL.Path)
		if fname == "" || fname == "/" {
			fname = "fish-cake.html"
		}
		log.Println("Serving static: " + self.templates + fname)
		http.ServeFile(res, req, self.templates+fname)
	}

	// Fix all routing
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", StaticHandler)
	r.HandleFunc("/favicon.ico", StaticHandler)
	r.HandleFunc("/ws", self.SocketHandler())
	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir(self.templates))))

	// Channel of updates
	update_clients := make(chan Ship)

	// Update all players on a time
	go func() {
		for _ = range ticker.C {
			for _, v := range self.players {
				v.Tick()
				update_clients <- v.ship
			}
		}
	}()

	// Simulation
	// Simulation
	// Simulation

	self.simulation.Start(33 * time.Millisecond)
	defer self.simulation.Stop()

	for i := -4.0; i <= 4.0; i++ {
		for j := -4.0; j <= 4.0; j++ {
			for k := -4.0; k <= 4.0; k++ {
				x := i
				y := j
				z := k
				self.simulation.Add(
					Vector3{x, y, z},
					Vector3{0.0, 0.0, 1.0},
				)
			}
		}
	}

	// Send swifts updates
	go func() {
		for m := range self.simulation.channel {
			for _, p := range self.players {
				p.SendMessage(m)
			}
		}
	}()

	//
	go func() {
		for s := range update_clients {
			w := World{
				nil,
				map[uint64]Ship{s.Id: s},
				-1.0,
				-1.0,
			}
			encoded, _ := json.Marshal(w)
			for _, v := range self.players {
				// v.connection.WriteMessage(1, []byte(encoded))
				_ = v
				_ = encoded
			}
		}
	}()

	log.Println("nevermind again -- bia next time")
	log.Println(STATIC_DIR)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

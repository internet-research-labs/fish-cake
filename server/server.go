package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"math"
	"net/http"
	"path"
	"time"
)

const (
	STATIC_DIR = "/static/"
	PORT       = "8080"
)

// GameHamndler
func GameHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../server/static/index.html"))
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, "ok")
}

// GameHandler
func SocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	s0 := Ship{
		Id:    0,
		Coord: WorldCoord{math.Pi / 4.0, 0.0},
	}

	s1 := Ship{
		Id:    1,
		Coord: WorldCoord{-math.Pi / 4.0, 0.0},
	}

	// Read loop
	go func() {
		for {
			t, m, e := conn.ReadMessage()
			if e != nil {
				return
			}
			log.Println(t, m, e)
			log.Println("you might ask yourself")
		}
	}()

	ticker := time.NewTicker(13 * time.Millisecond)

	for _ = range ticker.C {
		s0.Coord.Fi = math.Mod(s0.Coord.Fi+0.05, 2*math.Pi)
		s1.Coord.Fi = math.Mod(s1.Coord.Fi-0.05, 2*math.Pi)
		ships := make(map[uint64]Ship)
		ships[s0.Id] = s0
		ships[s1.Id] = s1
		res, err := json.Marshal(World{
			Ships:     ships,
			Radius:    10.0,
			Thickness: 1.0,
		})
		if err == nil {
			conn.WriteMessage(1, []byte(res))
		}
	}
}

// Server holds local zone game information
type Server struct {
	zones   []Zone
	players []uint
	bots    []uint
	world   World
}

// NewRandomServer returns a new server with a random world
func NewRandomServer(n int) Server {
	return Server{
		zones:   nil,
		players: nil,
		bots:    nil,
		world:   RandomWorld(n),
	}
}

// upgrader adds websocket support
var upgrader = websocket.Upgrader{}

// StaticHandler responds to requests by path
func StaticHandler(res http.ResponseWriter, req *http.Request) {
	fname := path.Base(req.URL.Path)
	http.ServeFile(res, req, "../server/static/"+fname)
}

// Listen for connections
func (self *Server) Listen(port int) {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", GameHandler)
	r.HandleFunc("/ws", SocketHandler)
	r.HandleFunc("/favicon.ico", StaticHandler)
	r.HandleFunc("/world", func(w http.ResponseWriter, r *http.Request) {
		fs, err := json.Marshal(self.world)
		if err != nil {
			w.Write([]byte("sicko mode"))
			return
		}
		w.Write([]byte(fs))
	})

	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("../server"+STATIC_DIR))))

	log.Println("nevermind again -- bia next time")
	log.Println(STATIC_DIR)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

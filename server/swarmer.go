package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type SwarmParams struct {
	Attraction float64 `json:"attraction"`
	Repulsion  float64 `json:"repulsion"`
	Alignment  float64 `json:"alignment"`
}

type ConnectionFamily struct {
	conns map[string]*websocket.Conn
	mutex sync.Mutex
	In    chan SwiftMap
}

func NewConnectionFamily() *ConnectionFamily {

	family := ConnectionFamily{
		conns: make(map[string]*websocket.Conn),
		In:    make(chan SwiftMap),
	}

	go func() {
		for update := range family.In {
			for id, conn := range family.conns {
				log.Println("UPDATING:", id)
				conn.WriteMessage(1, EncodeWireMessage("yupdate", update))
			}
		}
	}()

	return &family
}

func (family *ConnectionFamily) Add(id string, conn *websocket.Conn) {
	family.mutex.Lock()
	defer family.mutex.Unlock()
	family.conns[id] = conn
}

func (family *ConnectionFamily) Remove(id string) {
	family.mutex.Lock()
	defer family.mutex.Unlock()
	delete(family.conns, id)
}

// SocketHandler returns a handler function for gorilla that adds a new ship
func SwarmSocketHandler() func(http.ResponseWriter, *http.Request) {

	zone := NewSwiftZone(0.001062, 0.002555, 0.491596)

	log.Println("Adding swifts")
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			for k := 0; k < 8; k++ {
				zone.Add(
					RandomVector3(-2.0, 2.0),
					RandomVector3(-0.02, 0.02),
				)
			}
		}
	}

	counter := 0

	// Seed before calling random
	rand.Seed(time.Now().UTC().UnixNano())
	zone.Start(3 * time.Millisecond)

	family := NewConnectionFamily()

	go func() {
		for m := range zone.channel {
			family.In <- m
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {

		counter += 1
		id := fmt.Sprintf("%d", counter)

		// Upgrade to a websocket connection
		conn, _ := upgrader.Upgrade(w, r, nil)

		family.Add(id, conn)

		defer func() {
			log.Println("Closing connection")
			family.Remove(id)
			conn.Close()
		}()

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}
}

// Listen for connections
func SwarmListen(port int, templates string) {

	log.Println(fmt.Sprintf("Listening on %d with static files: %s", port, templates))

	// StaticHandler responds to requests by path
	// NOTE: Needs to be in closure because of self reference
	StaticHandler := func(res http.ResponseWriter, req *http.Request) {
		fname := path.Base(req.URL.Path)
		if fname == "" || fname == "/" {
			fname = "swarm.html"
		}
		log.Println("Serving static: " + templates + fname)
		http.ServeFile(res, req, templates+fname)
	}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", StaticHandler)
	r.HandleFunc("/favicon.ico", StaticHandler)
	r.HandleFunc("/swarm", SwarmSocketHandler())
	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir(templates))))

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

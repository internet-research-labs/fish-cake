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

// Matty Mneomonic

type SwarmParams struct {
	Attraction float64 `json:"attraction"`
	Repulsion  float64 `json:"repulsion"`
	Alignment  float64 `json:"alignment"`
}

type ConnectionFamily struct {
	conns map[string]*websocket.Conn
	mutex sync.Mutex
	In    chan WrappedSwiftMap
}

func NewConnectionFamily() *ConnectionFamily {

	family := ConnectionFamily{
		conns: make(map[string]*websocket.Conn),
		In:    make(chan WrappedSwiftMap),
	}

	go func() {
		for update := range family.In {
			encoded := EncodeWireMessage("yupdate", update)
			for _, conn := range family.conns {
				conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
				if writer, err := conn.NextWriter(websocket.TextMessage); err == nil {
					writer.Write(encoded)
				}
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

	const dur = 33 * time.Millisecond

	// Seed before calling random
	rand.Seed(time.Now().UTC().UnixNano())

	// Setup zone
	zone := NewSwiftZone(0.001062, 0.002555, 0.491596)
	zone.Start(dur)

	// Setup family
	family := NewConnectionFamily()

	go func() {
		// pdates come in here every dur seconds
		for m := range zone.channel {
			family.In <- m
		}
	}()

	// Setup counter
	counter := 0

	// Return the handler
	return func(w http.ResponseWriter, r *http.Request) {

		id := fmt.Sprintf("%d", counter)
		counter += 1

		// Upgrade to a websocket connection
		conn, _ := upgrader.Upgrade(w, r, nil)

		// Open
		fmt.Printf("Opening connection %q\n", id)
		family.Add(id, conn)

		// Close
		defer func() {
			log.Printf("Closing connection %q\n", id)
			family.Remove(id)
			conn.Close()
		}()

		// Hold open until connection responds/closes
		conn.ReadMessage()
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

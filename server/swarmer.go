package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"path"
	"time"
)

// SocketHandler returns a handler function for gorilla that adds a new ship
func SwarmSocketHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Upgrade to a websocket connection
		conn, _ := upgrader.Upgrade(w, r, nil)

		log.Println("Connecting")

		log.Println("Making your own personal world")
		zone := NewSwiftZone(0.001, 0.0014, 0.5)
		zone.Start(30 * time.Millisecond)

		// Seed before calling random
		rand.Seed(time.Now().UTC().UnixNano())

		log.Println("Adding swifts")
		for i := 0.0; i < 8.0; i++ {
			for j := 0.0; j < 8.0; j++ {
				for k := 0.0; k < 8.0; k++ {
					/*
						x := i - 4.0
						y := j - 4.0
						z := k - 4.0
					*/
					zone.Add(
						// Vector3{x, y, z},
						// Vector3{0.0, 0.0, -0.2},
						RandomVector3(-2.0, 2.0),
						RandomVector3(-0.02, 0.02),
					)
				}
			}
		}

		defer func() {
			log.Println("Closing connection")
			zone.Stop()
			conn.Close()
		}()

		go func() {
			for m := range zone.channel {
				conn.WriteMessage(1, EncodeWireMessage("yupdate", m))
			}
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
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

	// Fix all routing
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", StaticHandler)
	r.HandleFunc("/favicon.ico", StaticHandler)
	r.HandleFunc("/swarm", SwarmSocketHandler())
	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir(templates))))

	// ok
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

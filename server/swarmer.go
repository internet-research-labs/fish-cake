package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
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
		zone := NewSwiftZone()
		zone.Start(30 * time.Millisecond)

		log.Println("Adding swifts")
		for i := -4.0; i <= 4.0; i++ {
			for j := -4.0; j <= 4.0; j++ {
				for k := -4.0; k <= 4.0; k++ {
					x, y, z := i, j, k
					zone.Add(
						Vector3{x, y, z},
						Vector3{0.0, 0.0, 0.0},
					)
				}
			}
		}

		defer func() {
			log.Println("Closing connection")
			zone.Stop()
			conn.Close()
		}()

		for m := range zone.channel {
			conn.WriteMessage(1, EncodeWireMessage("yupdate", m))
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

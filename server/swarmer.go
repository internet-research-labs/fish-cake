package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"path"
	"time"
)

type SwarmParams struct {
	Attraction float64 `json:"attraction"`
	Repulsion  float64 `json:"repulsion"`
	Alignment  float64 `json:"alignment"`
}

// SocketHandler returns a handler function for gorilla that adds a new ship
func SwarmSocketHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(*r)

		// Seed before calling random
		rand.Seed(time.Now().UTC().UnixNano())

		// Upgrade to a websocket connection
		conn, _ := upgrader.Upgrade(w, r, nil)

		log.Println("Connecting")

		log.Println("Making your own personal world")
		zone := NewSwiftZone(
			Random(0.001, 0.003),
			Random(0.001, 0.003),
			Random(0.0, 1.0),
		)
		zone.Start(30 * time.Millisecond)

		log.Println("Adding swifts")
		for i := 0.0; i < 8.0; i++ {
			for j := 0.0; j < 8.0; j++ {
				for k := 0.0; k < 8.0; k++ {
					zone.Add(
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
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			} else {
				params := SwarmParams{}
				json.Unmarshal(msg, &params)
				zone.Attraction = params.Attraction
				zone.Repulsion = params.Repulsion
				zone.Alignment = params.Alignment
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

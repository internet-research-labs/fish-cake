// Handler returns a websocket http handler
package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
)

// Matty Mneomonic

type SwarmParams struct {
	Attraction float64 `json:"attraction"`
	Repulsion  float64 `json:"repulsion"`
	Alignment  float64 `json:"alignment"`
}

func StartZone(caster *Broadcaster) {
	const dur = 33 * time.Millisecond

	// Seed before calling random
	rand.Seed(time.Now().UTC().UnixNano())

	// Setup zone
	zone := NewSwiftZone(
		Random(0, 0.5),
		Random(0, 1.0),
		Random(0, 1.0),
	)

	zone.Start(dur)

	// Setup family
	go func() {
		// pdates come in here every dur seconds
		for m := range zone.channel {
			caster.In() <- EncodeWireMessage("yupdate", m)
		}
	}()
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

	// Static
	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir(templates))))

	// Generic broadcast handler
	broadcaster := NewBroadcaster()
	broadcaster.Start()
	r.HandleFunc("/swarm", broadcaster.Handler())

	StartZone(broadcaster)

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

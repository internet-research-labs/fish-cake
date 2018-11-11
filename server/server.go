package server

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"time"
)

const (
	STATIC_DIR = "/static/"
	PORT       = "8080"
)

// SocketHandler
func GameHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../server/static/index.html"))
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, "ok")
}

// GameHandler
func SocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

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

	ticker := time.NewTicker(1 * time.Second)

	for _ = range ticker.C {
		conn.WriteMessage(1, []byte("whatever"))
	}
}

// Things that the game needs to keep track of
type Server struct {
	zones   []Zone
	players []uint
	bots    []uint
}

var upgrader = websocket.Upgrader{}

// Listen for connections
func (self *Server) Listen(port int) {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", GameHandler)
	r.HandleFunc("/ws", SocketHandler)

	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("../server"+STATIC_DIR))))

	log.Println("nevermind again -- bia next time")
	log.Println(STATIC_DIR)
	log.Fatal(http.ListenAndServe(":8000", r))
}

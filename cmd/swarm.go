package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/internet-research-labs/fish-cake/server"
)

var (
	static *string = flag.String("static", "", "directory of static assets")
	port   *int    = flag.Int("port", 8080, "port to listen on")
)

func main() {

	flag.Parse()

	if static == nil || *static == "" {
		fmt.Println("Must set -static to a valid directory. Exiting nonzero")
		os.Exit(1)
	}

	server.SwarmListen(*port, *static)
}

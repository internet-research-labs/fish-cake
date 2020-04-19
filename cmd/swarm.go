package main

import (
	"flag"
	"github.com/internet-research-labs/fish-cake/server"
)

func main() {
	var static = flag.String("static", "", "directory of static assets")
	flag.Parse()
	server.SwarmListen(8080, *static)
}

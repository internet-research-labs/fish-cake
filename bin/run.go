package main

import (
	"github.com/internet-research-labs/fish-cake/server"
)

func main() {
	s := server.Server{}
	s.Listen(8080)
}

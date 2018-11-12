package main

import (
	"github.com/internet-research-labs/fish-cake/server"
)

func main() {
	s := server.NewRandomServer(10)
	s.Listen(8080)
}

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/internet-research-labs/fish-cake/server"
)

var (
	rows int = 25
	cols int = 100
)

func main() {

	ticker := time.NewTicker(1 * time.Second)

	game := server.NewGameOfLife(rows, cols)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	r.Float64()

	game.Add(1, 1)
	game.Add(1, 2)
	game.Add(2, 1)
	game.Add(2, 2)

	// fmt.Println(game.LivingCount(3, 3))
	// fmt.Println(game.LivingCount(3, 4))
	// fmt.Println(game.LivingCount(4, 3))
	// fmt.Println(game.LivingCount(4, 4))

	// Set initial state
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if r.Float64() < 0.1 {
				game.Add(i, j)
			}
		}
	}

	for _ = range ticker.C {
		fmt.Println(game)
		game.Tick()
	}

	ticker.Stop()
}

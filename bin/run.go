package main

import (
	"flag"
	"github.com/internet-research-labs/fish-cake/server"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
)

func main() {
	//
	var cpuprofile = flag.String("profile", "", "write cpu profile to `file`")

	// Command-line
	flag.Parse()

	// CPU profiling
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				log.Println("Trapped interrupt", sig)
				pprof.StopCPUProfile()
				os.Exit(0)
			}
		}()

		log.Println("Profiling!")
	}

	s := server.NewRandomServer(10)
	s.Listen(8080)
}

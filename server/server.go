// +build !client

package main

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/engine"
	"git.emersion.fr/saucisse-royale/miko.git/server/server"
	"git.emersion.fr/saucisse-royale/miko.git/server/timeserver"
	"log"
	"os"
	"os/signal"
)

// Miko server
func main() {
	address := ":9999"

	// Create server & engine
	log.Println("Creating server with address", address)
	srv := server.New(address)
	timeSrv := timeserver.New(":9998")
	e := engine.New(srv, timeSrv)
	ctx := e.Context()

	// Listen for SIGINT
	go (func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		for {
			<-c
			log.Println("Stopping server")
			e.Stop()
			os.Exit(0)
		}
	})()

	// Generate a new terrain
	ctx.Terrain.Generate()

	// Start server & engine

	go (func() {
		err := srv.Listen()
		if err != nil {
			log.Fatal("Server error:", err)
		}
	})()

	go (func() {
		err := timeSrv.Listen()
		if err != nil {
			log.Fatal("Time server error:", err)
		}
	})()

	e.Start()
}

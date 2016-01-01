// +build !client

package main

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/engine"
	"git.emersion.fr/saucisse-royale/miko.git/server/server"
	"git.emersion.fr/saucisse-royale/miko.git/server/timeserver"
	"os"
	"os/signal"
)

// Miko server
func main() {
	// Create server & engine
	srv := server.New(":9999")
	timeSrv := timeserver.New(":9998")
	e := engine.New(srv, timeSrv)
	ctx := e.Context()

	// Listen for SIGINT
	go (func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		for {
			<-c
			e.Stop()
			os.Exit(0)
		}
	})()

	// Generate a new terrain
	ctx.Terrain.Generate()

	// Start server & engine
	go srv.Listen()
	go timeSrv.Listen()
	e.Start()
}

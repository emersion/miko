// +build !client

package main

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/engine"
	"git.emersion.fr/saucisse-royale/miko.git/server/server"
)

// Miko server
func main() {
	srv := server.New(":9999")

	e := engine.New(srv)
	ctx := e.Context()

	// Generate a new terrain
	ctx.Terrain.Generate()

	// Start the server & the engine
	go srv.Listen()
	e.Start()
}

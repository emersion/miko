// +build !client

package main

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/engine"
	"git.emersion.fr/saucisse-royale/miko.git/server/server"
)

// Miko server
func main() {
	e := engine.New()
	ctx := e.Context()

	// Generate a new terrain
	ctx.Terrain.Generate()

	// Start the server
	srv := server.New(":9999", ctx)
	go srv.Listen()

	e.Start()
}

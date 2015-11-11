// +build !client

package main

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/auth"
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/server"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

// Miko server
func main() {
	// Create and generate a new terrain
	trn := terrain.New()
	trn.Generate()

	// Create a new context, with corresponding services
	ctx := message.NewServerContext()
	ctx.Auth = auth.NewService()
	ctx.Entity = entity.NewService()
	ctx.Terrain = trn
	ctx.Clock = clock.NewService()

	go ctx.Entity.Animate(ctx.Terrain, ctx.Clock)

	// Start the server
	srv := server.New(":9999", ctx)
	srv.Listen()
}

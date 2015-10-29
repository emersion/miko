package main

import(
	"git.emersion.fr/saucisse-royale/miko/server/auth"
	"git.emersion.fr/saucisse-royale/miko/server/entity"
	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/server"
	"git.emersion.fr/saucisse-royale/miko/server/terrain"
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

	go ctx.Entity.Animate(trn)

	// Start the server
	srv := server.New(":9999", ctx)
	srv.Listen()
}

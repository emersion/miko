package main

import(
	"git.emersion.fr/saucisse-royale/miko/server/auth"
	"git.emersion.fr/saucisse-royale/miko/server/entity"
	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/server"
	"git.emersion.fr/saucisse-royale/miko/server/terrain"
)

func main() {
	ctx := &message.Context{
		Type: message.ServerContext,
		Auth: auth.NewService(),
		Entity: entity.NewService(),
		Terrain: terrain.New(),
	}

	srv := server.New(":9999", ctx)
	srv.Listen()
}

package main

import (
	"git.emersion.fr/saucisse-royale/miko/server/auth"
	"git.emersion.fr/saucisse-royale/miko/server/entity"
	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/terrain"
	"git.emersion.fr/saucisse-royale/miko/server/browser/server"
)

func main() {
	address := ":9998"

	trn := terrain.New()
	trn.Generate()

	ctx := &message.Context{
		Type: message.ServerContext,
		Auth: auth.NewService(),
		Entity: entity.NewService(),
		Terrain: trn,
	}

	srv := server.New(address, ctx)
	srv.Listen()
}

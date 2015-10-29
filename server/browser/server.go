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

	ctx := message.NewServerContext()
	ctx.Auth = auth.NewService()
	ctx.Entity = entity.NewService()
	ctx.Terrain = trn

	go ctx.Entity.Animate(trn)

	srv := server.New(address, ctx)
	srv.Listen()
}

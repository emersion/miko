package main

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/auth"
	"git.emersion.fr/saucisse-royale/miko.git/server/browser/server"
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
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

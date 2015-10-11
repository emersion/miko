package main

import (
	"git.emersion.fr/saucisse-royale/miko/server/auth"
	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/terrain"
	"git.emersion.fr/saucisse-royale/miko/server/browser/server"
)

func main() {
	address := ":9998"

	ctx := &message.Context{
		Auth: auth.NewService(),
		Terrain: terrain.New(),
	}

	srv := server.New(address, ctx)
	srv.Listen()
}

package main

import(
	"./auth"
	"./message"
	"./server"
	"./terrain"
)

func main() {
	ctx := &message.Context{
		Auth: auth.NewService(),
		Terrain: terrain.New(),
	}

	srv := server.New(":9999", ctx)
	srv.Listen()
}

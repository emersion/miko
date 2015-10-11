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
	message.SetContext(ctx)

	srv := server.New(":9999")
	srv.Listen()
}

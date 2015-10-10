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

	srv := server.New("127.0.0.1:9999")
	srv.Listen()
}

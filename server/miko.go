package main

import(
	"./message"
	"./server"
)

func main() {
	ctx := &message.Context{
		Auth: server.NewAuthService(),
	}
	message.SetContext(ctx)

	srv := server.New("127.0.0.1:9999")
	srv.Listen()
}

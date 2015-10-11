package main

import (
	"net/http"
	"golang.org/x/net/websocket"

	"../auth"
	"../message"
	"../terrain"
)

func WsServer(ws *websocket.Conn) {
	clientIO := &message.IO{
		Reader: ws,
		Writer: ws,
		BroadcastWriter: ws,
		Id: 0,
	}

	message.Listen(clientIO)
}

func main() {
	ctx := &message.Context{
		Auth: auth.NewService(),
		Terrain: terrain.New(),
	}
	message.SetContext(ctx)

	http.Handle("/socket", websocket.Handler(WsServer))
	http.Handle("/", http.FileServer(http.Dir("public")))
	err := http.ListenAndServe(":9998", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

package main

import (
	"log"
	"net/http"
	"golang.org/x/net/websocket"

	"../auth"
	"../message"
	"../message/handler"
	"../terrain"
)

var hdlr handler.Handler

func WsServer(ws *websocket.Conn) {
	clientIO := &message.IO{
		Reader: ws,
		Writer: ws,
		BroadcastWriter: ws,
		Id: 0,
	}

	hdlr.Listen(clientIO)
}

func main() {
	address := ":9998"

	ctx := &message.Context{
		Auth: auth.NewService(),
		Terrain: terrain.New(),
	}
	hdlr := handler.New(ctx)

	log.Println("Creating HTTP server with address", address)
	http.Handle("/socket", websocket.Handler(WsServer))
	http.Handle("/", http.FileServer(http.Dir("public")))
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

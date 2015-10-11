package main

import (
	"encoding/binary"
	"net/http"
	"log"
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

	var msg_type message.Type
	for {
		err := binary.Read(ws, binary.BigEndian, &msg_type)
		if err != nil {
			ws.Close()
			log.Println("binary.Read failed:", err)
			return
		}
		message.Handle(msg_type, clientIO)
	}
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

package main

import(
	//"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"

	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/message/handler"
	"git.emersion.fr/saucisse-royale/miko/server/message/builder"
)

func main() {
	ctx := &message.Context{}
	hdlr := handler.New(ctx)

	c, err := websocket.Dial("ws://localhost:9998/socket")
	if err != nil {
		panic("Dial: " + err.Error())
	}

	clientIO := &message.IO{
		Reader: c,
		Writer: c,
	}

	go hdlr.Listen(clientIO)

	err = builder.SendPing(clientIO.Writer)
	if err != nil {
		panic("Write: " + err.Error())
	}

	/*err = c.Close()
	if err != nil {
		panic("Close: " + err.Error())
	}*/
}

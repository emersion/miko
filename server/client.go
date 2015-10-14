package main

import(
	"bufio"
	"net"
	//"log"

	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/message/handler"
	"git.emersion.fr/saucisse-royale/miko/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko/server/terrain"
)

func main() {
	trn := terrain.New()
	ctx := &message.Context{
		Type: message.ClientContext,
		Terrain: trn,
	}
	hdlr := handler.New(ctx)

	c, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		panic("Dial: " + err.Error())
	}

	reader := bufio.NewReader(c)
	clientIO := &message.IO{
		Reader: reader,
		Writer: c,
	}

	go hdlr.Listen(clientIO)

	/*err = builder.SendPing(clientIO.Writer)
	if err != nil {
		panic("SendPing: " + err.Error())
	}*/

	err = builder.SendLogin(clientIO.Writer, "root", "root")
	if err != nil {
		panic("SendLogin: " + err.Error())
	}

	for {}

	/*err = builder.SendChatSend(clientIO.Writer, "Hello World!")
	if err != nil {
		panic("SendChatSend: " + err.Error())
	}*/

	/*err = c.Close()
	if err != nil {
		panic("Close: " + err.Error())
	}*/
}

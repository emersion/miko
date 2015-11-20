// +build client

package main

import (
	"bufio"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"

	"git.emersion.fr/saucisse-royale/miko.git/server/browser/client"
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/handler"
)

func main() {
	trn := client.NewTerrain(js.Global.Get("document").Call("getElementById", "terrain"))
	ent := client.NewEntityService(js.Global.Get("document").Call("getElementById", "entities"))

	ctx := message.NewClientContext()
	ctx.Terrain = trn
	ctx.Entity = ent
	ctx.Clock = clock.NewService()
	ctx.Me = &message.Session{Username: "root"}

	hdlr := handler.New(ctx)

	host := js.Global.Get("window").Get("location").Get("host").String()
	c, err := websocket.Dial("ws://" + host + "/socket")
	if err != nil {
		panic("Dial: " + err.Error())
	}

	reader := bufio.NewReader(c)
	clientIO := message.NewIO(0, reader, c, nil)

	go hdlr.Listen(clientIO)

	engine := client.NewEngine(ctx, c)

	err = builder.SendVersion(clientIO)
	if err != nil {
		panic("SendVersion: " + err.Error())
	}

	/*err = builder.SendPing(clientIO)
	if err != nil {
		panic("SendPing: " + err.Error())
	}*/

	// TODO: wait for login response
	err = builder.SendLogin(clientIO, "root", "root")
	if err != nil {
		panic("SendLogin: " + err.Error())
	}

	engine.Start()

	/*err = builder.SendChatSend(clientIO, "Hello World!")
	if err != nil {
		panic("SendChatSend: " + err.Error())
	}*/

	/*err = c.Close()
	if err != nil {
		panic("Close: " + err.Error())
	}*/
}

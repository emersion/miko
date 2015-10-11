package main

import(
	//"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"

	"git.emersion.fr/saucisse-royale/miko/message"
)

func main() {
	c, err := websocket.Dial("ws://localhost:9998/socket")
	if err != nil {
		panic("Dial: " + err.Error())
	}

	err = message.SendPing(c)
	if err != nil {
		panic("Write: " + err.Error())
	}

	/*err = c.Close()
	if err != nil {
		panic("Close: " + err.Error())
	}*/
}

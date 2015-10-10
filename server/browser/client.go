package main

import(
	//"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
)

func main() {
	c, err := websocket.Dial("ws://localhost/socket")
	if err != nil {
		panic("Dial: " + err.Error())
	}

	_, err = c.Write([]byte("Hello!"))
	if err != nil {
		panic("Write: " + err.Error())
	}

	err = c.Close()
	if err != nil {
		panic("Close: " + err.Error())
	}
}

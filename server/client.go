package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"

	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/crypto"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/handler"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

// A very basic TCP client, for testing purposes
func main() {
	trn := terrain.New()

	ctx := message.NewClientContext()
	ctx.Terrain = trn
	ctx.Clock = clock.NewService()

	hdlr := handler.New(ctx)

	tlsConfig, err := crypto.GetClientTlsConfig()
	if err != nil {
		log.Println("WARN: could not load TLS config")
	}

	var c net.Conn
	if tlsConfig != nil {
		c, err = tls.Dial("tcp", "127.0.0.1:9999", tlsConfig)
	} else {
		c, err = net.Dial("tcp", "127.0.0.1:9999")
	}

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

	for {
	}

	/*err = builder.SendChatSend(clientIO.Writer, "Hello World!")
	if err != nil {
		panic("SendChatSend: " + err.Error())
	}*/

	/*err = c.Close()
	if err != nil {
		panic("Close: " + err.Error())
	}*/
}

package main

import "./server"

func main() {
	srv := server.New("127.0.0.1:9999")
	srv.Listen()
}

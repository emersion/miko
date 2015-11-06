all: server client

client:
	gopherjs build client.go -o public/client.js
server:
	go build -o miko server.go
start:
	go run server.go
dev:
	make all && ./miko

.PHONY: all client server start

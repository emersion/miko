all: server client

client:
	gopherjs build -o public/client.js client.go
server:
	go build -o miko
start:
	go run
dev:
	make all && ./miko

.PHONY: all client server start

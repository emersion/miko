all: server
server:
	go build -o miko server.go
client:
	go build -o miko-client client.go
start:
	go run server.go
start-client:
	go run client.go

.PHONY: all server client start start-client

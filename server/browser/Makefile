all: client server

client:
	# TODO: cd $GOPATH/src/git.emersion.fr/saucisse-royale; ln -s /home/XXX/projects/miko/server/ miko
	gopherjs build client.go -o public/client.js
server:
	go build -o miko server.go
start:
	go run server.go

.PHONY: all client server start

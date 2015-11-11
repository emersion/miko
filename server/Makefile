all: server
server:
	go build -o miko
client:
	go build -o miko-client -tags client
start:
	go run server.go
start-client:
	go run client.go
crypto:
	openssl genrsa -out crypto/server.key 2048
	openssl req -new -x509 -key crypto/server.key -out crypto/server.pem -days 3650
test:
	go test ./auth ./clock ./entity

.PHONY: all server client start start-client crypto

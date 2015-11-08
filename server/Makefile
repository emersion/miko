all: server
server:
	go build -o miko server.go
client:
	go build -o miko-client client.go
start:
	go run server.go
start-client:
	go run client.go
crypto:
	openssl genrsa -out crypto/server.key 2048
	openssl req -new -x509 -key crypto/server.key -out crypto/server.pem -days 3650
test:
	go test git.emersion.fr/saucisse-royale/miko.git/server/auth
	go test git.emersion.fr/saucisse-royale/miko.git/server/entity

.PHONY: all server client start start-client crypto

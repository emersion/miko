ifeq ($(OS),Windows_NT)
	SUDO=
else
	SUDO=sudo
endif

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
	go test ./...
docker:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s' -o miko-docker
	$(SUDO) docker build -t miko .
start-docker:
	$(SUDO) docker run -it miko

clean:
	rm -f miko miko-client miko-docker

.PHONY: all server client start start-client crypto docker start-docker

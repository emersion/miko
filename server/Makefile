all:
	go build server.go -o miko
start:
	go run server.go

browser-client:
	gopherjs build browser/client.go -o browser/client.js
browser-server:
	go build browser/server.go -o browser/server
browser-start:
	go run browser/server.go

.PHONY: all start browser-client browser-server browser-start

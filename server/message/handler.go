package message

import (
	"errors"
)

type Handler func(*IO) error

var handlers = map[Type]Handler{}

func main() {
	handlers[Ping] = func(io *IO) error {
		SendPingResp(io.Writer)
		return nil
	}
}

func Handle(t Type, io *IO) error {
	if val, ok := handlers[t]; ok {
		return val(io);
	} else {
		return errors.New("Unknown message type")
	}
}

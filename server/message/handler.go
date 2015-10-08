package message

import (
	"errors"
)

type Handler func(*IO) error

var handlers = map[Type]Handler{}

func main() {
	handlers[Ping] = func(io *IO) error {
		return SendPingResp(io.Writer)
	}

	handlers[Exit] = func(io *IO) error {
		return io.Writer.Close()
	}
}

func Handle(t Type, io *IO) error {
	if val, ok := handlers[t]; ok {
		return val(io);
	} else {
		return errors.New("Unknown message type")
	}
}

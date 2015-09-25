package message

import (
	"io"
	"bufio"
	"errors"
)

type Handler func(bufio.Reader, io.Writer) error

var handlers = map[Type]Handler{}

func main() {
	handlers[Ping] = func(reader bufio.Reader, writer io.Writer) error {
		SendPingResp(writer)
		return nil
	}
}

func Handle(t Type, reader bufio.Reader, writer io.Writer) error {
	if val, ok := handlers[t]; ok {
		return val(reader, writer);
	} else {
		return errors.New("Unknown message type")
	}
}

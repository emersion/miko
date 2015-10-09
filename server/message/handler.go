package message

import (
	"io"
	"errors"
	"encoding/binary"
)

func read(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.BigEndian, data)
}

func readString(r io.Reader) string {
	var len uint8
	err := read(r, &len)
	if err != nil {
		return ""
	}

	out := make([]byte, len)
	n, err := r.Read(out)
	if n < int(len) || err != nil {
		return ""
	}

	return string(out)
}

var ctx Context

type Handler func(*IO) error

var handlers = map[Type]Handler{
	Types["ping"]: func(io *IO) error {
		return SendPingResp(io.Writer)
	},
	Types["exit"]: func(io *IO) error {
		return io.Writer.Close()
	},
	Types["login"]: func(io *IO) error {
		username := readString(io.Reader)
		password := readString(io.Reader)

		code := ctx.Auth.Login(io, username, password)
		if err := SendLoginResp(io.Writer, code); err != nil {
			return err
		}

		if code == LoginResponseCodes["ok"] {
			return SendPlayerJoined(io.BroadcastWriter, username)
		} else {
			return nil
		}
	},
	Types["chatsend"]: func(io *IO) error {
		//msg := readString(io.Reader)
		return nil
	},
}

func Handle(t Type, io *IO) error {
	if val, ok := handlers[t]; ok {
		return val(io);
	} else {
		return errors.New("Unknown message type")
	}
}

func SetContext(c *Context) {
	ctx = *c
}

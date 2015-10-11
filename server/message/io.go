package message

import (
	"io"
	"encoding/binary"
)

type IO struct {
	Reader io.Reader
	Writer io.WriteCloser
	BroadcastWriter io.Writer
	Id int
}

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

func write(w io.Writer, data interface{}) error {
	return binary.Write(w, binary.BigEndian, data)
}

func writeString(w io.Writer, data string) error {
	if err := write(w, uint8(len(data))); err != nil {
		return err;
	}
	return write(w, data)
}

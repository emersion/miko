package handler

import (
	"encoding/binary"
	"errors"
	"io"
)

func read(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.BigEndian, data)
}

func readString(r io.Reader) (str string, err error) {
	var len uint16
	err = read(r, &len)
	if err != nil {
		return
	}

	out := make([]byte, len)
	n, err := r.Read(out)
	if err != nil {
		return
	}
	if n < int(len) {
		err = errors.New("Could not read enough data while reading string")
		return
	}

	str = string(out)
	return
}

func Read(r io.Reader, data ...interface{}) error {
	for _, item := range data {
		if str, ok := item.(*string); ok {
			out, err := readString(r)
			if err != nil {
				return err
			}
			*str = out
		} else {
			err := read(r, item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

package handler

import (
	"io"
	"log"
	"encoding/binary"
)

func read(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.BigEndian, data)
}

func readString(r io.Reader) string {
	var len uint8
	err := read(r, &len)
	if err != nil {
		log.Println("WARN: cannot read string length")
		return ""
	}

	out := make([]byte, len)
	n, err := r.Read(out)
	if n < int(len) || err != nil {
		log.Println("WARN: cannot read string, expected length:", len, "but could only read:", n)
		return ""
	}

	return string(out)
}

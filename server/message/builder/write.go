package builder

import (
	"encoding/binary"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"io"
	"log"
)

func write(w io.Writer, data interface{}) error {
	if t, ok := data.(message.Type); ok {
		log.Println("Sent:", message.GetTypeName(t))
	}

	return binary.Write(w, binary.BigEndian, data)
}

func writeString(w io.Writer, data string) error {
	if err := write(w, uint16(len(data))); err != nil {
		return err
	}

	_, err := w.Write([]byte(data))
	if err != nil {
		return err
	}

	return nil
}

func writeAll(w io.Writer, data []interface{}) error {
	for _, item := range data {
		var err error
		switch item.(type) {
		case string:
			err = writeString(w, item.(string))
		default:
			err = write(w, item)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func lock(w io.Writer) {
	if client, ok := w.(*message.IO); ok {
		client.Locker.Lock()
	}
}

func unlock(w io.Writer) {
	if client, ok := w.(*message.IO); ok {
		client.Locker.Unlock()
	}
}

func send(w io.Writer, data interface{}) error {
	lock(w)
	defer unlock(w)

	return write(w, data)
}

func sendAll(w io.Writer, data []interface{}) error {
	lock(w)
	defer unlock(w)

	return writeAll(w, data)
}

package builder

import (
	"encoding/binary"
	"io"
	"sync"
	//"log"
	//"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

func write(w io.Writer, data interface{}) error {
	/*if t, ok := data.(message.Type); ok {
		log.Println("Sent:", message.GetTypeName(t))
	}*/

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

func Write(w io.Writer, data ...interface{}) error {
	for _, item := range data {
		if str, ok := item.(string); ok {
			err := writeString(w, str)
			if err != nil {
				return err
			}
		} else {
			err := write(w, item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func lock(w io.Writer) {
	if locker, ok := w.(sync.Locker); ok {
		locker.Lock()
	} else {
		panic("Called lock() on a non-locker variable")
	}
}

func unlock(w io.Writer) {
	if locker, ok := w.(sync.Locker); ok {
		locker.Unlock()
	} else {
		panic("Called unlock() on a non-locker variable")
	}
}

func send(w io.Writer, data ...interface{}) error {
	lock(w)
	defer unlock(w)

	return Write(w, data...)
}

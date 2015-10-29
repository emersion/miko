package builder

import (
	"io"
	"encoding/binary"
)

func write(w io.Writer, data interface{}) error {
	return binary.Write(w, binary.BigEndian, data)
}

func writeString(w io.Writer, data string) error {
	if err := write(w, uint16(len(data))); err != nil {
		return err;
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
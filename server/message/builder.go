package message

import(
	"io"
	"encoding/binary"
)

func write(w io.Writer, data interface{}) error {
	return binary.Write(w, binary.BigEndian, data)
}

func writeString(w io.Writer, data string) error {
	if err := write(w, uint8(len(data))); err != nil {
		return err;
	}
	return write(w, data)
}

func SendPingResp(w io.Writer) error {
	return write(w, GetRespType(Types["ping"]))
}

func SendLoginResp(w io.Writer, code string) error {
	if err := write(w, GetRespType(Types["login"])); err != nil {
		return err
	}
	return write(w, LoginResponseCode[code])
}

func SendPlayerJoined(w io.Writer, username string) error {
	if err := write(w, Types["playermeta"]); err != nil {
		return err
	}
	if err := write(w, MetaActionCode["playerjoined"]); err != nil {
		return err
	}
	return writeString(w, username)
}

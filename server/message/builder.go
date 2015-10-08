package message

import(
	"io"
	"encoding/binary"
)

func send(writer io.Writer, data interface{}) error {
	return binary.Write(writer, binary.BigEndian, data)
}

func SendPingResp(writer io.Writer) error {
	return send(writer, GetRespType(Ping))
}

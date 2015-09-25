package message

import(
	"io"
	"encoding/binary"
)

func SendPingResp(writer io.Writer) {
	binary.Write(writer, binary.BigEndian, GetRespType(Ping))
}

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
	return write(w, Types["pong"])
}

func SendLoginResp(w io.Writer, code LoginResponseCode) error {
	if err := write(w, Types["loginresponse"]); err != nil {
		return err
	}
	return write(w, code)
}

func SendRegisterResp(w io.Writer, code RegisterResponseCode) error {
	if err := write(w, Types["registerresponse"]); err != nil {
		return err
	}
	return write(w, code)
}

func SendPlayerMeta(w io.Writer, code MetaActionCode, username string) error {
	if err := write(w, Types["playermeta"]); err != nil {
		return err
	}
	if err := write(w, code); err != nil {
		return err
	}
	return writeString(w, username)
}
func SendPlayerJoined(w io.Writer, username string) error {
	return SendPlayerMeta(w, MetaActionCodes["playerjoined"], username)
}
func SendPlayerLeft(w io.Writer, username string) error {
	return SendPlayerMeta(w, MetaActionCodes["playerleft"], username)
}

func SendTerrainUpdate(w io.Writer, block *Block) error {
	// Get the point type the most used
	var defaultType PointType
	defaultTypeCount := -1
	typesStats := make(map[PointType]int)
	for _, column := range block.Points {
		for _, ptType := range column {
			if _, present := typesStats[ptType]; !present {
				typesStats[ptType] = 1
			} else {
				typesStats[ptType]++
			}

			if defaultTypeCount < typesStats[ptType] {
				defaultTypeCount = typesStats[ptType]
				defaultType = ptType
			}
		}
	}

	// Send the response
	if err := write(w, Types["terrainupdate"]); err != nil {
		return err
	}
	if err := write(w, block.X); err != nil {
		return err
	}
	if err := write(w, block.Y); err != nil {
		return err
	}
	if err := write(w, defaultType); err != nil {
		return err
	}
	if err := write(w, uint16(block.Size() - defaultTypeCount)); err != nil {
		return err
	}

	for x, column := range block.Points {
		for y, ptType := range column {
			if ptType == defaultType {
				continue
			}

			if err := write(w, PointCoord(x)); err != nil {
				return err
			}
			if err := write(w, PointCoord(y)); err != nil {
				return err
			}
			if err := write(w, ptType); err != nil {
				return err
			}
		}
	}

	return nil
}

func SendChatReceive(w io.Writer, username string, msg string) error {
	if err := write(w, Types["chatreceive"]); err != nil {
		return err
	}
	if err := writeString(w, username); err != nil {
		return err
	}
	return writeString(w, msg)
}

package builder

import (
	"io"
	".."
)

func SendLoginResp(w io.Writer, code message.LoginResponseCode) error {
	if err := write(w, message.Types["login_response"]); err != nil {
		return err
	}
	return write(w, code)
}

func SendRegisterResp(w io.Writer, code message.RegisterResponseCode) error {
	if err := write(w, message.Types["register_response"]); err != nil {
		return err
	}
	return write(w, code)
}

func SendPlayerMeta(w io.Writer, code message.MetaActionCode, username string) error {
	if err := write(w, message.Types["player_meta"]); err != nil {
		return err
	}
	if err := write(w, code); err != nil {
		return err
	}
	return writeString(w, username)
}
func SendPlayerJoined(w io.Writer, username string) error {
	return SendPlayerMeta(w, message.MetaActionCodes["player_joined"], username)
}
func SendPlayerLeft(w io.Writer, username string) error {
	return SendPlayerMeta(w, message.MetaActionCodes["player_left"], username)
}

func SendTerrainUpdate(w io.Writer, block *message.Block) error {
	// Get the point type the most used
	var defaultType message.PointType
	defaultTypeCount := -1
	typesStats := make(map[message.PointType]int)
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
	if err := write(w, message.Types["terrain_update"]); err != nil {
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

			if err := write(w, message.PointCoord(x)); err != nil {
				return err
			}
			if err := write(w, message.PointCoord(y)); err != nil {
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
	if err := write(w, message.Types["chat_receive"]); err != nil {
		return err
	}
	if err := writeString(w, username); err != nil {
		return err
	}
	return writeString(w, msg)
}

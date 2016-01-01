package message_test

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"testing"
	"time"
)

func TestTimeToTimestamp(test *testing.T) {
	t := time.Unix(789697858, 137000)
	expected := message.Timestamp(789697858137)
	timestamp := message.TimeToTimestamp(t)

	if timestamp != expected {
		test.Error("Expected", expected, "but got", timestamp)
	}
}

func TestTimestampToTime(test *testing.T) {
	timestamp := message.Timestamp(807656789876568)
	expected := time.Unix(807656789876, 568000)
	t := message.TimestampToTime(timestamp)

	if t != expected {
		test.Error("Expected", expected, "but got", t)
	}
}

package game

import (
	"io"

	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/handler"
)

type Config struct {
	message.ConfigBase
	DefaultPlayerSpeed  float32
	PlayerBallCooldown  Cooldown
	DefaultBallSpeed    float32
	DefaultBallLifespan TicksLeft
}

func (c *Config) WriteTo(w io.Writer) (n int64, err error) {
	err = builder.Write(w, c.MaxRollbackTicks, c.DefaultPlayerSpeed,
		c.PlayerBallCooldown, c.DefaultBallSpeed, c.DefaultBallLifespan)
	return
}

func (c *Config) ReadFrom(r io.Reader) (n int64, err error) {
	err = handler.Read(r, &c.MaxRollbackTicks, &c.DefaultPlayerSpeed,
		&c.PlayerBallCooldown, &c.DefaultBallSpeed, &c.DefaultBallLifespan)
	return
}

func DefaultConfig() *Config {
	return &Config{
		ConfigBase: message.ConfigBase{
			MaxRollbackTicks: uint16(message.MaxRewind),
		},
		DefaultPlayerSpeed:  7,
		PlayerBallCooldown:  20,
		DefaultBallSpeed:    9,
		DefaultBallLifespan: 100,
	}
}

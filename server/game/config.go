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
	err = builder.Write(w, c.MaxRollbackTicks, c.TimeServerPort, c.DefaultPlayerSpeed,
		c.PlayerBallCooldown, c.DefaultBallSpeed, c.DefaultBallLifespan)
	return
}

func (c *Config) ReadFrom(r io.Reader) (n int64, err error) {
	err = handler.Read(r, &c.MaxRollbackTicks, &c.TimeServerPort, &c.DefaultPlayerSpeed,
		&c.PlayerBallCooldown, &c.DefaultBallSpeed, &c.DefaultBallLifespan)
	return
}

func DefaultConfig() *Config {
	config := &Config{}
	config.MaxRollbackTicks = uint16(message.MaxRewind)
	config.DefaultPlayerSpeed = 7
	config.PlayerBallCooldown = 20
	config.DefaultBallSpeed = 9
	config.DefaultBallLifespan = 100
	return config
}

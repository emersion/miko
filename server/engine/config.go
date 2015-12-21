package engine

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

type Config struct {
	MaxRollbackTicks    uint16
	DefaultPlayerSpeed  float32
	PlayerBallCooldown  uint16
	DefaultBallSpeed    float32
	DefaultBallLifespan uint16
}

func (c *Config) Export() []interface{} {
	return []interface{}{
		c.MaxRollbackTicks,
		c.DefaultPlayerSpeed,
		c.PlayerBallCooldown,
		c.DefaultBallSpeed,
		c.DefaultBallLifespan,
	}
}

func (c *Config) Import(data []interface{}) {
	c.MaxRollbackTicks = data[0].(uint16)
	c.DefaultPlayerSpeed = data[1].(float32)
	c.PlayerBallCooldown = data[2].(uint16)
	c.DefaultBallSpeed = data[3].(float32)
	c.DefaultBallLifespan = data[4].(uint16)
}

func DefaultConfig() *Config {
	return &Config{
		MaxRollbackTicks:    uint16(message.MaxRewind),
		DefaultPlayerSpeed:  7,
		PlayerBallCooldown:  20,
		DefaultBallSpeed:    9,
		DefaultBallLifespan: 100,
	}
}

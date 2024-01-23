package core

import "github.com/nanderv/traincontrol-prototype/internal/bridge"

type Configurator func(*Core) error

func WithFakeBridge() func(c *Core) error {
	return func(c *Core) error {
		c.CommandBridge = bridge.NewFakeBridge(&MessageAdapter{c: c})

		return nil
	}
}

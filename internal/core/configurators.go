package core

import "github.com/nanderv/traincontrol-prototype/internal/bridge"

type Configurator func(*Core) error

func WithFakeBridge() func(c *Core) error {
	return func(c *Core) error {
		c.commandBridge = bridge.NewFakeBridge(&MessageAdapter{c: c})

		return nil
	}
}

func WithTrackSwitch(id byte) func(c *Core) error {
	return func(c *Core) error {
		c.trackSwitches = append(c.trackSwitches, TrackSwitch{number: id})
		return nil
	}
}

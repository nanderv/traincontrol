package core

import "github.com/nanderv/traincontrol-prototype/internal/core/domain/layout"

type Configurator func(*Core) error

func WithTrackSwitch(id byte) func(c *Core) error {
	return func(c *Core) error {
		c.layout.TrackSwitches = append(c.layout.TrackSwitches, layout.TrackSwitch{Number: id})
		return nil
	}
}

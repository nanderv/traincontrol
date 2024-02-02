package traintracks

import "github.com/nanderv/traincontrol-prototype/internal/traintracks/domain/layout"

type Configurator func(*TrackService) error

func WithTrackSwitch(id byte) func(c *TrackService) error {
	return func(c *TrackService) error {
		c.layout.TrackSwitches = append(c.layout.TrackSwitches, layout.TrackSwitch{Number: id})
		return nil
	}
}

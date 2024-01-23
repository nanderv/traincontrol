package core

import "github.com/nanderv/traincontrol-prototype/internal/bridge"

type Configurator func(Core) (Core, error)

func WithFakeBridge() func(c Core) (Core, error) {
	return func(c Core) (Core, error) {
		c.CommandBridge, c.CommandBridgeReturnChannel = bridge.NewFakeBridge()

		return c, nil
	}
}

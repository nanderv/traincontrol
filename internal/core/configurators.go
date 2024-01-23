package core

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/types"
)

type Configurator func(Core) (Core, error)

func WithFakeBridge() func(c Core) (Core, error) {
	return func(c Core) (Core, error) {
		cc := make(chan types.Msg, 10)
		bridge := bridge.FakeBridge{
			Result: &cc,
		}
		c.CommandBridge = &bridge
		return c, nil
	}
}

package core

import (
	"context"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/types"
)

type Core struct {
	CommandBridge              CommandBridge
	CommandBridgeReturnChannel *chan types.Msg
}

func NewCore(configurator ...Configurator) (Core, error) {
	c := Core{}

	for _, config := range configurator {
		var err error
		c, err = config(c)
		if err != nil {
			return Core{}, err
		}
	}
	return c, nil
}

func (c *Core) SetSwitch(switchID byte, direction bool) {
	c.CommandBridge.Send(types.SetSwitch{SwitchID: switchID, Direction: direction}.ToBridgeMsg())
}

func (c *Core) EventHandler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-(*c.CommandBridgeReturnChannel):
			fmt.Println("OUT", msg)
		}
	}
}

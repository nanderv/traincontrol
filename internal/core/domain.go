package core

import (
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/types"
)

type Core struct {
	CommandBridge        CommandBridge
	MessageReturnChannel *chan bridge.Msg
}

func NewCore(configurator ...Configurator) (*Core, error) {
	c := Core{}
	ch := make(chan bridge.Msg, 10)
	c.MessageReturnChannel = &ch
	for _, config := range configurator {
		var err error
		err = config(&c)
		if err != nil {
			return &Core{}, err
		}
	}
	return &c, nil
}

func (c *Core) SetSwitch(switchID byte, direction bool) error {
	return c.CommandBridge.Send(types.SetSwitch{SwitchID: switchID, Direction: direction}.ToBridgeMsg())
}
func (c *Core) HandleSwitchSet(msg types.SetSwitchResult) {
	fmt.Println(msg.String())
}

package core

import (
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/types"
)

type MessageAdapter struct {
	c *Core
}

func (m *MessageAdapter) SendReturnMessage(msg bridge.Msg) error {
	fmt.Println("OUT", msg)

	switch msg.Type {
	case 3:
		vv := types.SetSwitchResult{SetSwitch: types.NewSetSwitch(msg.Val[0], msg.Val[1] == 1)}

		m.c.HandleSwitchSet(vv)
	}
	return nil
}

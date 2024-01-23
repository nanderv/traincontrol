package core

import (
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/types"
)

func (c *Core) handleSwitchSet(msg types.Msg) {
	v := types.SetSwitchResult{SetSwitch: types.SetSwitch{SwitchID: msg.Val[0], Direction: msg.Val[1] == 1}}.String()

	fmt.Println(v)
}

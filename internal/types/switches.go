package types

import (
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
)

type SetSwitch struct {
	SwitchID  byte
	Direction bool
}

func (s SetSwitch) ToBridgeMsg() bridge.Msg {
	var d bridge.Msg
	d.Type = 2
	d.Val[0] = s.SwitchID
	if s.Direction {
		d.Val[1] = 1
	} else {
		d.Val[1] = 0
	}
	return d
}

func (s SetSwitch) String() string {
	v := "left"
	if s.Direction {
		v = "right"
	}
	return fmt.Sprintf("Switch %v set to direction %s", s.SwitchID, v)
}

type SetSwitchResult struct {
	SetSwitch
}

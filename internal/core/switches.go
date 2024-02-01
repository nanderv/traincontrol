package core

import (
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
)

func NewSetSwitch(switchId byte, direction bool) SetSwitch {
	return SetSwitch{
		switchID:  switchId,
		direction: direction,
	}
}

type SetSwitch struct {
	switchID  byte
	direction bool
}

func (s SetSwitch) ToBridgeMsg() domain.Msg {
	var d domain.Msg
	d.Type = 2
	d.Val[0] = s.switchID
	if s.direction {
		d.Val[1] = 1
	} else {
		d.Val[1] = 0
	}
	return d
}

func (s SetSwitch) String() string {
	v := "left"
	if s.direction {
		v = "right"
	}
	return fmt.Sprintf("Switch %v set to Direction %s", s.switchID, v)
}

type SetSwitchResult struct {
	SetSwitch
}

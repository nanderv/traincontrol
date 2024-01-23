package types

type SetSwitch struct {
	SwitchID  byte
	Direction bool
}

func (s SetSwitch) ToBridgeMsg() Msg {
	var d Msg
	d.Type = 2
	d.Val[0] = s.SwitchID
	if s.Direction {
		d.Val[1] = 0
	} else {
		d.Val[1] = 1
	}
	return d
}

type SetSwitchResult struct {
	SetSwitch
}

func (s SetSwitchResult) ToBridgeMsg() Msg {
	d := s.SetSwitch.ToBridgeMsg()
	d.Type = d.Type + 1
	return d
}

package core

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"log/slog"
)

type MessageAdapter struct {
	c *Core
}

func (m *MessageAdapter) SendReturnMessage(msg bridge.Msg) error {
	slog.Info("INCOMING", "Data", msg)

	switch msg.Type {
	case 3:
		vv := SetSwitchResult{SetSwitch: NewSetSwitch(msg.Val[0], msg.Val[1] == 1)}

		m.c.SetSwitchEvent(vv)
	}
	return nil
}

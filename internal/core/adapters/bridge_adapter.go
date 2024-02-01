package adapters

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/core"
	"log/slog"
)

type MessageAdapter struct {
	c *core.Core
	r core.SendCommand
}

func (m *MessageAdapter) Receive(msg domain.Msg) error {
	slog.Info("INCOMING", "Data", msg)

	switch msg.Type {
	case 3:
		vv := core.SetSwitchResult{SetSwitch: core.NewSetSwitch(msg.Val[0], msg.Val[1] == 1)}

		m.c.SetSwitchEvent(vv)
	}
	return nil
}
func (m *MessageAdapter) Send(msg domain.Msg) error {
	return m.r.Send(msg)
}
func NewMessageAdapter(c *core.Core, b *bridge.SerialBridge) *MessageAdapter {
	m := MessageAdapter{c: c, r: b}
	c.AddCommandBridge(&m)
	b.AddReceiver(&m)
	return &m
}

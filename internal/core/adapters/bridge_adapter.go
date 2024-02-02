package adapters

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/core"
	"log/slog"
)

type MessageAdapter struct {
	core   *core.Core
	sender core.MessageSender
}

// Receive a message from a layout
func (ma *MessageAdapter) Receive(msg domain.Msg) error {
	slog.Info("INCOMING", "Data", msg)

	switch msg.Type {
	case 3:
		c := core.SetSwitchResult{SetSwitch: core.NewSetSwitch(msg.Val[0], msg.Val[1] == 1)}
		ma.core.SetSwitchEvent(c)
	}
	return nil
}

// Send a message towards a layout
func (ma *MessageAdapter) Send(msg domain.Msg) error {
	return ma.sender.Send(msg)
}

func NewMessageAdapter(c *core.Core, b *bridge.SerialBridge) *MessageAdapter {
	m := MessageAdapter{core: c, sender: b}
	c.AddCommandBridge(&m)
	b.AddReceiver(&m)
	return &m
}

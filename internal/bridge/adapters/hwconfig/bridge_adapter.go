package hwconfig

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain/codes"
	"github.com/nanderv/traincontrol-prototype/internal/hwconfig"
	"github.com/nanderv/traincontrol-prototype/internal/hwconfig/domain/node"
	"github.com/nanderv/traincontrol-prototype/internal/hwconfig/domain/subcodes"
	"log/slog"
)

type MessageAdapter struct {
	core   *hwconfig.HwConfigurator
	sender bridge.Bridge
}

// Receive a message from a layout
func (ma *MessageAdapter) Receive(msg domain.Msg) error {
	if msg.Type != codes.HW {
		return nil
	}
	switch msg.Val[0] {
	case 1:
		ma.core.HandleNodeAnnounce([3]byte{msg.Val[1], msg.Val[2], msg.Val[3]}, msg.Val[4])
	}
	return nil
}

// Send a message towards a layout
func (ma *MessageAdapter) Send(msg domain.Msg) error {
	return ma.sender.Send(msg)
}
func (ma *MessageAdapter) SendNodeInfoUpdate(node node.Node) error {
	slog.Info("Sending node info for node", "node", node)
	msg := domain.Msg{
		Type: 0,
		Val:  [6]byte{subcodes.AckAnnounce, node.Mac[0], node.Mac[1], node.Mac[2], node.Addr, 166},
	}
	return ma.sender.Send(msg)
}
func NewMessageAdapter(hwConfig *hwconfig.HwConfigurator, b bridge.Bridge) *MessageAdapter {
	m := MessageAdapter{core: hwConfig, sender: b}
	hwConfig.SetBridgeAdapter(&m)
	b.AddReceiver(&m)
	return &m
}

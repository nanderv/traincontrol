package hwconfig

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain/codes"
	"github.com/nanderv/traincontrol-prototype/internal/hwconfig"
	"github.com/nanderv/traincontrol-prototype/internal/hwconfig/domain/node"
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
		slog.Info("NODE REGISTERED", "Mac", [3]byte{msg.Val[1], msg.Val[2], msg.Val[3]}, "AddrRequested", msg.Val[4])
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
	return nil
}
func NewMessageAdapter(hwConfig *hwconfig.HwConfigurator, b bridge.Bridge) *MessageAdapter {
	m := MessageAdapter{core: hwConfig, sender: b}
	hwConfig.SetBridgeAdapter(&m)
	b.AddReceiver(&m)
	return &m
}

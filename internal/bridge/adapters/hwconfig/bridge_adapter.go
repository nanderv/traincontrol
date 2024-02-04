package hwconfig

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain/codes"
	"github.com/nanderv/traincontrol-prototype/internal/hwconfig"
	"github.com/nanderv/traincontrol-prototype/internal/hwconfig/domain/node"
	"github.com/nanderv/traincontrol-prototype/internal/hwconfig/domain/subcodes"
	"log/slog"
	"time"
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
func (ma *MessageAdapter) SendNodeInfoUpdate(node node.Node, state byte) error {
	slog.Info("Sending node info for node", "node", node)
	msg := domain.Msg{
		Type: 0,
		Val:  [6]byte{subcodes.AckAnnounce, node.Mac[0], node.Mac[1], node.Mac[2], node.Addr, state},
	}

	return ma.sender.SendWithResponseChecksAndRetries(msg, func(d domain.Msg) bool {
		return d.Type == 0 &&
			d.Val[0] == subcodes.Announce &&
			d.Val[1] == msg.Val[1] &&
			d.Val[2] == msg.Val[2] &&
			d.Val[3] == msg.Val[3] &&
			d.Val[4] == msg.Val[4] &&
			d.Val[5] == msg.Val[5]
	}, 7*time.Second, 10)
}

func (ma *MessageAdapter) SendEepromWrite(nodeAddr byte, commandSlot byte, typ byte, data [2]byte) error {
	msg := domain.Msg{
		Type: 0,
		Val:  [6]byte{subcodes.EEPROM_WRITE, nodeAddr, commandSlot, typ, data[0], data[1]},
	}
	return ma.sender.SendWithResponseChecksAndRetries(msg, func(d domain.Msg) bool {
		return d.Type == 0 &&
			d.Val[0] == subcodes.EEPROM_WRITE_RETURN &&
			d.Val[1] == nodeAddr &&
			d.Val[2] == commandSlot &&
			d.Val[3] == typ &&
			d.Val[4] == data[0] &&
			d.Val[5] == data[1]
	}, 100*time.Millisecond, 10)
}

func NewMessageAdapter(hwConfig *hwconfig.HwConfigurator, b bridge.Bridge) *MessageAdapter {
	m := MessageAdapter{core: hwConfig, sender: b}
	hwConfig.SetBridgeAdapter(&m)
	b.AddReceiver(&m)
	return &m
}

package traintracks

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge/adapters"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain/codes"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks/domain/commands"
	"log/slog"
	"time"
)

type MessageAdapter struct {
	core *traintracks.TrackService

	sender    bridgeSender[domain.Msg]
	listeners map[*chan domain.Msg]struct{}
}

func NewMessageAdapter(svc *traintracks.TrackService, bridge Bridge) *MessageAdapter {
	m := MessageAdapter{core: svc, sender: bridge, listeners: make(map[*chan domain.Msg]struct{})}
	svc.SetLayoutSender(&m)
	bridge.AddReceiver(&m)
	return &m
}

// Receive a message from a layout
func (ma *MessageAdapter) Receive(msg domain.Msg) error {
	slog.Info("INCOMING", "Data", msg)
	ma.sendToListners(msg)

	return ma.handleReceivedMessage(msg)
}

func (ma *MessageAdapter) handleReceivedMessage(msg domain.Msg) error {
	switch msg.Type {
	case codes.HW:
		return nil
	case codes.SwitchResult:
		c := commands.SetSwitchResult{SetSwitch: commands.NewSetSwitch(msg.Val[0], msg.Val[1] == 1)}
		ma.core.UpdateSwitchState(c)
	}
	return nil
}

func (ma *MessageAdapter) SetSwitchDirection(switchID byte, direction bool) error {
	retriesRemaining := 10

	cha := ma.addListener()
	defer ma.removeListener(cha)

	msg := commands.NewSetSwitch(switchID, direction)

	resultChecker := func(m domain.Msg) bool {
		return m.Type == 3 && m.Val[0] == switchID && (m.Val[1] == 1) == direction
	}

	return adapters.SendMessageWithConfirmationAndRetries(ma.sender.Send, cha, resultChecker, msg.ToBridgeMsg(), retriesRemaining, 500*time.Millisecond)
}

func (ma *MessageAdapter) addListener() *chan domain.Msg {
	ch := make(chan domain.Msg)
	ma.listeners[&ch] = struct{}{}
	return &ch
}

func (ma *MessageAdapter) removeListener(ch *chan domain.Msg) {
	delete(ma.listeners, ch)
	return
}
func (ma *MessageAdapter) sendToListners(msg domain.Msg) {
	for lnr, _ := range ma.listeners {
		*lnr <- msg
	}
}

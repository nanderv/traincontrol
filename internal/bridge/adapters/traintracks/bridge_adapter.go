package traintracks

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/adapters"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain/codes"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks/domain/commands"
	"log/slog"
	"sync"
	"time"
)

type MessageAdapter struct {
	sync.Mutex

	trackService *traintracks.TrackService
	sender       bridge.Bridge

	listeners map[*chan domain.Msg]struct{}
}

func NewMessageAdapter(svc *traintracks.TrackService, bridge bridge.Bridge) *MessageAdapter {
	m := MessageAdapter{trackService: svc, sender: bridge, listeners: make(map[*chan domain.Msg]struct{})}
	svc.SetLayoutSender(&m)
	bridge.AddReceiver(&m)
	return &m
}

// Receive a message from a layout
func (adapt *MessageAdapter) Receive(msg domain.Msg) error {
	slog.Info("INCOMING", "Data", msg)
	adapt.sendToListeners(msg)

	return adapt.handleReceivedMessage(msg)
}

func (adapt *MessageAdapter) handleReceivedMessage(msg domain.Msg) error {
	switch msg.Type {
	case codes.HW:
		return nil
	case codes.SwitchResult:
		c := commands.SetSwitchResult{SetSwitch: commands.NewSetSwitch(msg.Val[0], msg.Val[1] == 1)}
		adapt.trackService.UpdateSwitchState(c)
	}
	return nil
}

func (adapt *MessageAdapter) SetSwitchDirection(switchID byte, direction bool) error {
	listener := adapt.addListener()
	defer adapt.removeListener(listener)

	msg := commands.NewSetSwitch(switchID, direction)

	sender := adapters.Sender{
		Bridge: adapt.sender,
		ResultChecker: func(m domain.Msg) bool {
			return m.Type == codes.SwitchResult && m.Val[0] == switchID && (m.Val[1] == 1) == direction
		},
		ListenChannel: listener,
	}

	requestTimeout := 300 * time.Second
	retries := 10
	return sender.SendMessageWithConfirmationAndRetries(msg.ToBridgeMsg(), requestTimeout, retries)
}

func (adapt *MessageAdapter) addListener() *chan domain.Msg {
	adapt.Lock()
	defer adapt.Unlock()

	ch := make(chan domain.Msg)
	adapt.listeners[&ch] = struct{}{}
	return &ch
}

func (adapt *MessageAdapter) removeListener(ch *chan domain.Msg) {
	adapt.Lock()
	defer adapt.Unlock()

	delete(adapt.listeners, ch)
}

func (adapt *MessageAdapter) sendToListeners(msg domain.Msg) {
	adapt.Lock()
	defer adapt.Unlock()
	for lnr, _ := range adapt.listeners {
		*lnr <- msg
	}
}

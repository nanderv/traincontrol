package traintracks

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain/codes"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks/domain/commands"
	"time"
)

type MessageAdapter struct {
	trackService *traintracks.TrackService
	sender       bridge.Bridge
}

func NewMessageAdapter(svc *traintracks.TrackService, bridge bridge.Bridge) *MessageAdapter {
	m := MessageAdapter{trackService: svc, sender: bridge}
	svc.SetLayoutSender(&m)
	bridge.AddReceiver(&m)
	return &m
}

// Receive a message from a layout
func (adapt *MessageAdapter) Receive(msg domain.Msg) error {
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
	msg := commands.NewSetSwitch(switchID, direction)

	checker := func(m domain.Msg) bool {
		return m.Type == codes.SwitchResult && m.Val[0] == switchID && (m.Val[1] == 1) == direction
	}

	requestTimeout := 300 * time.Millisecond
	retries := 10
	return adapt.sender.SendWithResponseChecksAndRetries(msg.ToBridgeMsg(), checker, requestTimeout, retries)
}

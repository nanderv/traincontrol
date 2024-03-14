package adapters

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks"
	domain2 "github.com/nanderv/traincontrol-prototype/internal/traintracks/domain"
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
	case domain.HW:
		return nil
	case domain.SwitchResult:
		sw, err := adapt.trackService.Layout.GetSwitchFromHWIDs(domain.Mac{msg.Val[0], msg.Val[1], msg.Val[2]}, msg.Val[3], msg.Val[4])
		if err != nil {
			return err
		}

		return adapt.trackService.UpdateSwitchState(sw, msg.Val[5] != 0)
	}
	return nil
}

func (adapt *MessageAdapter) SetSwitchDirection(t *domain2.TrackSwitch, dir bool) error {
	msg := t.SetDirectionCMD(dir)

	responseChecker := func(m domain.Msg) bool {
		return m.Type == domain.SwitchResult && m.Val == msg.Val
	}

	requestTimeout := 600 * time.Millisecond
	retries := 3
	return adapt.sender.SendWithResponseChecksAndRetries(msg, responseChecker, requestTimeout, retries)
}

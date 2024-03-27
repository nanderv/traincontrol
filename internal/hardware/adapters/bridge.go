package adapters

import (
	"context"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/hardware"
	domain2 "github.com/nanderv/traincontrol-prototype/internal/hardware/domain"
	"log/slog"
	"time"
)

type Bridge interface {
	AddListener(ctx context.Context) *chan domain.Msg
	Send(domain.Msg) error
	SendWithResponseChecksAndRetries(msg domain.Msg, checker func(msg domain.Msg) bool, timeout time.Duration, retries int) error
}

type Adapter struct {
	trackService *hardware.TrackService
	sender       Bridge
}

func NewAdapter(svc *hardware.TrackService, bridge Bridge) *Adapter {
	m := Adapter{trackService: svc, sender: bridge}
	svc.SetLayoutSender(&m)

	cha := bridge.AddListener(context.Background())
	go m.handle(cha)
	return &m
}

func (a *Adapter) handle(cha *chan domain.Msg) {
	for {
		select {
		case msg := <-*cha:
			err := a.Receive(msg)
			if err != nil {
				slog.Error("Error", "err", err)
			}
		}
	}
}

// Receive a message from a layout
func (a *Adapter) Receive(msg domain.Msg) error {
	switch msg.Type {
	case domain.HW:
		return nil
	case domain.SwitchResult:
		sw, err := a.trackService.Layout.GetSwitchFromHWIDs(domain.Mac{msg.Val[0], msg.Val[1], msg.Val[2]}, msg.Val[3], msg.Val[4])
		if err != nil {
			return err
		}

		return a.trackService.UpdateSwitchState(sw, msg.Val[5] != 0)
	}
	return nil
}

func (a *Adapter) SetSwitchDirection(t *domain2.TrackSwitch, dir bool) error {
	msg := t.SetDirectionCMD(dir)

	responseChecker := func(m domain.Msg) bool {
		return m.Type == domain.SwitchResult && m.Val == msg.Val
	}

	requestTimeout := 1000 * time.Millisecond
	retries := 10
	return a.sender.SendWithResponseChecksAndRetries(msg, responseChecker, requestTimeout, retries)
}

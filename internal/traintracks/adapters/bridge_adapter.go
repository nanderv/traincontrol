package adapters

import (
	"errors"
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
	for lnr, _ := range ma.listeners {
		*lnr <- msg
	}
	slog.Info("INCOMING", "Data", msg)

	switch msg.Type {
	case codes.HW:
		return nil
	case codes.SwitchResult:
		c := commands.SetSwitchResult{SetSwitch: commands.NewSetSwitch(msg.Val[0], msg.Val[1] == 1)}
		ma.core.SetSwitchEvent(c)
	}
	return nil
}

func (ma *MessageAdapter) SetSwitchDirection(switchID byte, direction bool) error {
	retriesRemaining := 10

	cha := ma.addListener()
	defer ma.removeListener(cha)

	msg := commands.NewSetSwitch(switchID, direction)

	for retriesRemaining > 0 {
		err := ma.sender.Send(msg.ToBridgeMsg())
		if err != nil {
			return err
		}
		retriesRemaining--
		select {
		case resultMsg := <-*cha:
			if resultMsg.Type == 3 && resultMsg.Val[0] == switchID && (resultMsg.Val[1] == 1) == direction {
				slog.Info("Done direction", "message", resultMsg)
				return nil
			} else {
				// Correctly arrived messages that are not the right one don't count towards retry counter
				slog.Debug("Message received, but irrelevant", "message", resultMsg)
				retriesRemaining += 1
			}
		case <-time.After(500 * time.Millisecond):
			slog.Warn("timeout while sending ", "message", msg.ToBridgeMsg())
			break
		}
	}

	return errors.New("out of retries")
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

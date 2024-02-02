package adapters

import (
	"errors"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
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

func (ma *MessageAdapter) addListener() *chan domain.Msg {
	ch := make(chan domain.Msg)
	ma.listeners[&ch] = struct{}{}
	return &ch
}

func (ma *MessageAdapter) removeListener(ch *chan domain.Msg) {
	delete(ma.listeners, ch)
	return
}

// Receive a message from a layout
func (ma *MessageAdapter) Receive(msg domain.Msg) error {
	for r, _ := range ma.listeners {
		*r <- msg
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
func (ma *MessageAdapter) setSwitchDirection(switchID byte, direction bool, retriesRemaining int) error {
	ch := ma.addListener()
	defer ma.removeListener(ch)

	msg := commands.NewSetSwitch(switchID, direction)

	for retriesRemaining > 0 {
		err := ma.sender.Send(msg.ToBridgeMsg())
		if err != nil {
			return err
		}
		retriesRemaining--
		select {
		case msg := <-*ch:
			if msg.Type == 3 && msg.Val[0] == switchID && (msg.Val[1] == 1) == direction {
				fmt.Println("Done direction", msg)

				return nil
			} else {
				// Correctly arrived messages that are not the right one don't count towards retry counter
				retriesRemaining += 1
			}
		case <-time.After(100 * time.Millisecond):
			fmt.Println("timeout ", retriesRemaining)
			break
		}
	}

	return errors.New("out of retries")
}
func (ma *MessageAdapter) SetSwitchDirection(switchID byte, direction bool) error {
	fmt.Println("SS")
	return ma.setSwitchDirection(switchID, direction, 10)
}

type Bridge interface {
	AddReceiver(bridge.MessageReceiver)
	Send(domain.Msg) error
}

func NewMessageAdapter(c *traintracks.TrackService, b Bridge) *MessageAdapter {
	m := MessageAdapter{core: c, sender: b, listeners: make(map[*chan domain.Msg]struct{})}
	c.SetLayoutSender(&m)
	b.AddReceiver(&m)
	return &m
}

package core

import (
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
)

type Core struct {
	commandBridge        CommandBridge
	messageReturnChannel *chan bridge.Msg
	notifyChangeChannels []*chan struct{}

	trackSwitches []TrackSwitch
}

func NewCore(configurator ...Configurator) (*Core, error) {
	c := Core{}
	ch := make(chan bridge.Msg, 10)
	c.messageReturnChannel = &ch
	c.trackSwitches = make([]TrackSwitch, 0)
	c.notifyChangeChannels = make([]*chan struct{}, 0)
	for _, config := range configurator {
		var err error
		err = config(&c)
		if err != nil {
			return &Core{}, err
		}
	}
	return &c, nil
}

func (c *Core) SetSwitchAction(switchID byte, direction bool) error {
	var found bool
	for _, sw := range c.trackSwitches {
		if sw.number == switchID {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("switch with id %v not found", switchID)
	}
	return c.commandBridge.Send(NewSetSwitch(switchID, direction).ToBridgeMsg())
}

func (c *Core) SetSwitchEvent(msg SetSwitchResult) {
	for _, sw := range c.trackSwitches {
		if sw.number == msg.SetSwitch.switchID {
			sw.direction = msg.SetSwitch.direction
		}
	}
	c.notify()
	fmt.Println(msg.String())
}

type TrackSwitch struct {
	number    byte
	direction bool
}

func (c *Core) notify() {
	for _, ch := range c.notifyChangeChannels {
		select {
		case *ch <- struct{}{}:
			return
		default:
			return
		}
	}
}

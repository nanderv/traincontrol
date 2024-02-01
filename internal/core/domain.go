package core

import (
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
)

type Layout struct {
	TrackSwitches []TrackSwitch
}
type TrackSwitch struct {
	Number    byte
	Direction bool
}

type Core struct {
	commandBridge        CommandBridge
	messageReturnChannel *chan bridge.Msg
	notifyChangeChannels []*chan Layout
	layout               Layout
}

func (c *Core) AddNewReturnChannel() *chan Layout {
	ch := make(chan Layout)
	c.notifyChangeChannels = append(c.notifyChangeChannels, &ch)
	return &ch
}
func NewCore(configurator ...Configurator) (*Core, error) {
	c := Core{}
	ch := make(chan bridge.Msg, 10)
	c.messageReturnChannel = &ch
	c.layout.TrackSwitches = make([]TrackSwitch, 0)
	c.notifyChangeChannels = make([]*chan Layout, 0)
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
	for _, sw := range c.layout.TrackSwitches {
		if sw.Number == switchID {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("switch with id %v not found", switchID)
	}
	return c.commandBridge.Send(NewSetSwitch(switchID, direction).ToBridgeMsg())
}

func (c *Core) SetSwitchEvent(msg SetSwitchResult) {
	for i, sw := range c.layout.TrackSwitches {
		if sw.Number == msg.SetSwitch.switchID {
			c.layout.TrackSwitches[i].Direction = msg.SetSwitch.direction
		}
	}
	c.notify()
	fmt.Println(msg.String())
}

func (c *Core) notify() {
	fmt.Println("NOTIFY")
	for _, ch := range c.notifyChangeChannels {
		fmt.Println(ch)
		select {
		case *ch <- c.layout:
			fmt.Println(c.layout)
			return
		}
	}
}

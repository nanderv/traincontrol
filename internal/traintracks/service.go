package traintracks

import (
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks/domain/commands"
	layout2 "github.com/nanderv/traincontrol-prototype/internal/traintracks/domain/layout"
)

type TrackService struct {
	layoutBridge         LayoutSender
	notifyChangeChannels []*chan layout2.Layout
	layout               layout2.Layout
}

func (svc *TrackService) AddNewReturnChannel() *chan layout2.Layout {
	ch := make(chan layout2.Layout)
	svc.notifyChangeChannels = append(svc.notifyChangeChannels, &ch)
	return &ch
}

func (svc *TrackService) SetLayoutSender(cc LayoutSender) {
	svc.layoutBridge = cc
	return
}

func NewCore(configurator ...Configurator) (*TrackService, error) {
	c := TrackService{}
	c.layout.TrackSwitches = make([]layout2.TrackSwitch, 0)
	c.notifyChangeChannels = make([]*chan layout2.Layout, 0)
	for _, config := range configurator {
		var err error
		err = config(&c)
		if err != nil {
			return &TrackService{}, err
		}
	}
	return &c, nil
}

func (svc *TrackService) SetSwitchAction(switchID byte, direction bool) error {
	var found bool
	for _, sw := range svc.layout.TrackSwitches {
		if sw.Number == switchID {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("switch with id %v not found", switchID)
	}

	return svc.layoutBridge.SetSwitchDirection(switchID, direction)
}

func (svc *TrackService) UpdateSwitchState(msg commands.SetSwitchResult) {
	for i, sw := range svc.layout.TrackSwitches {
		if msg.SetSwitch.IsSwitch(sw.Number) {
			svc.layout.TrackSwitches[i].Direction = msg.SetSwitch.GetDirection()
		}
	}
	svc.notify()
}

func (svc *TrackService) notify() {
	for _, ch := range svc.notifyChangeChannels {
		*ch <- svc.layout
	}
}

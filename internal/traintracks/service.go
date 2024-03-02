package traintracks

import (
	domain2 "github.com/nanderv/traincontrol-prototype/internal/traintracks/domain"
)

type TrackService struct {
	layoutBridge         Sender
	notifyChangeChannels []*chan domain2.Layout
	Layout               domain2.Layout
}

func (svc *TrackService) AddNewReturnChannel() *chan domain2.Layout {
	ch := make(chan domain2.Layout)
	svc.notifyChangeChannels = append(svc.notifyChangeChannels, &ch)
	return &ch
}

func (svc *TrackService) SetLayoutSender(cc Sender) {
	svc.layoutBridge = cc
	return
}

func NewTrackService(lay domain2.Layout) (*TrackService, error) {
	c := TrackService{}
	c.Layout = lay

	return &c, nil
}

func (svc *TrackService) SetSwitchDirection(switchID string, direction bool) error {
	sw, err := svc.Layout.GetSwitch(switchID)

	if err != nil {
		return err
	}

	return svc.layoutBridge.SetSwitchDirection(sw, direction)
}

func (svc *TrackService) UpdateSwitchState(sw *domain2.TrackSwitch, direction bool) error {
	sw.UpdateDirection(direction)
	svc.notify()
	return nil
}

func (svc *TrackService) notify() {
	for _, ch := range svc.notifyChangeChannels {
		*ch <- svc.Layout
	}
}

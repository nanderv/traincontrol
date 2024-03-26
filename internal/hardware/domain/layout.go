package domain

import (
	"errors"
	bridgeDomain "github.com/nanderv/traincontrol-prototype/internal/serialbridge/domain"
	"log/slog"
)

type HardwareState struct {
	TrackSwitches []*TrackSwitch
}

func NewHardwareState() HardwareState {
	l := HardwareState{
		TrackSwitches: make([]*TrackSwitch, 0),
	}
	return l
}
func (l *HardwareState) WithTrackSwitch(t TrackSwitch) {
	l.TrackSwitches = append(l.TrackSwitches, &t)
}

type TrackSwitch struct {
	//
	Mac bridgeDomain.Mac

	PortID    byte
	LeftPin   byte
	RightPin  byte
	Direction bool
	Name      string
}

func (l *HardwareState) GetSwitch(id string) (*TrackSwitch, error) {
	for _, ts := range l.TrackSwitches {
		if ts.Name == id {
			return ts, nil
		}
	}

	return nil, errors.New("could not find switch for return")

}
func (l *HardwareState) GetSwitchFromHWIDs(mac bridgeDomain.Mac, portID byte, pinID byte) (*TrackSwitch, error) {
	for _, sw := range l.TrackSwitches {
		if sw.Mac == mac && sw.PortID == portID && (sw.LeftPin == pinID || sw.RightPin == pinID) {
			return sw, nil
		}
	}
	return nil, errors.New("could not find switch for return")
}

func (t *TrackSwitch) UpdateDirection(dir bool) {
	slog.Debug("Updating direciton to", "Direction", dir, "Switch", t)
	t.Direction = dir
}

func (t *TrackSwitch) SetDirectionCMD(direction bool) bridgeDomain.Msg {
	pin := t.LeftPin
	dir := byte(0)
	slog.Debug("Direction when creating msg: ", "direction", direction)
	if direction {
		pin = t.RightPin
		dir = 1
	}
	return bridgeDomain.Msg{
		Type: bridgeDomain.SwitchSet,
		Val: [6]byte{
			t.Mac[0],
			t.Mac[1],
			t.Mac[2],
			t.PortID,
			pin,
			dir,
		},
	}
}

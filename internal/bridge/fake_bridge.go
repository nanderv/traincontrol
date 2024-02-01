package bridge

import (
	"fmt"
)

type Returner interface {
	SendReturnMessage(Msg) error
}

// The SerialBridge is responsible for translating commands towards things the railway can understand
type FakeBridge struct {
	Returner Returner
}

func (f *FakeBridge) Send(m Msg) error {
	fmt.Println("IN", m)
	r := m
	if m.Type == 2 {
		r.Type = 3
	}
	err := f.Returner.SendReturnMessage(r)
	if err != nil {
		return err
	}
	return nil
}

func NewFakeBridge(cc Returner) *FakeBridge {
	bridge := FakeBridge{
		Returner: cc,
	}
	return &bridge
}
func (f *FakeBridge) Handle() {
	return
}

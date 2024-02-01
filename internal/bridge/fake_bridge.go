package bridge

import (
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
)

// The SerialBridge is responsible for translating commands towards things the railway can understand
type FakeBridge struct {
	Returner Receiver
}

func (f *FakeBridge) Send(m domain.Msg) error {
	fmt.Println("IN", m)
	r := m
	if m.Type == 2 {
		r.Type = 3
	}
	err := f.Returner.Receive(r)
	if err != nil {
		return err
	}
	return nil
}

func NewFakeBridge(cc Receiver) *FakeBridge {
	bridge := FakeBridge{
		Returner: cc,
	}
	return &bridge
}
func (f *FakeBridge) Handle() {
	return
}

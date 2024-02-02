package bridge

import (
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
)

// The SerialBridge is responsible for translating commands towards things the railway can understand
type FakeBridge struct {
	Returner MessageReceiver
}

func (f *FakeBridge) AddReceiver(r MessageReceiver) {
	f.Returner = r
}
func (f *FakeBridge) Send(m domain.Msg) error {
	fmt.Println("IN", m)
	r := m
	if m.Type == 2 {
		r.Type = 3
	}
	go f.Returner.Receive(r)

	return nil
}

func NewFakeBridge() *FakeBridge {
	bridge := FakeBridge{}
	return &bridge
}
func (f *FakeBridge) Handle() {
	return
}

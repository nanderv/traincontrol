package bridge

import (
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/types"
)

// The FakeBridge is responsible for translating commands towards things the railway can understand
type FakeBridge struct {
	Result *chan types.Msg
}

func (f *FakeBridge) Send(m types.Msg) {
	fmt.Println("IN", m)
	r := m
	if m.Type == 2 {
		r.Type = 3
	}
	*f.Result <- r
}

func NewFakeBridge(cc *chan types.Msg) *FakeBridge {
	bridge := FakeBridge{
		Result: cc,
	}
	return &bridge
}

package bridge

import "github.com/nanderv/traincontrol-prototype/internal/types"

// The FakeBridge is responsible for translating commands towards things the railway can understand
type FakeBridge struct {
	result *chan types.Msg
}

func (f *FakeBridge) Send(m types.Msg) {
	r := m
	if m.Type == 1 {
		r.Type = 2
	}
	*f.result <- r
}

package bridge

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"log/slog"
)

// The SerialBridge is responsible for translating commands towards things the railway can understand
type FakeBridge struct {
	Returner []MessageReceiver
}

func (f *FakeBridge) AddReceiver(r MessageReceiver) {
	f.Returner = append(f.Returner, r)
}

func (f *FakeBridge) Send(m domain.Msg) error {
	slog.Info("OUTBOUND", "message", m)
	msg := m
	if m.Type == 2 {
		msg.Type = 3
	}

	go func() {
		for _, r := range f.Returner {
			err := r.Receive(msg)
			if err != nil {
				slog.Error("incorrect message", err)
			}
		}
	}()

	return nil
}

func NewFakeBridge() *FakeBridge {
	bridge := FakeBridge{}
	return &bridge
}
func (f *FakeBridge) Handle() {
	return
}

package serialbridge

import (
	"github.com/nanderv/traincontrol-prototype/internal/serialbridge/domain"
	"log/slog"
	"time"
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
	switch m.Type {
	case domain.SwitchSet:
		msg.Type = domain.SwitchResult
	case domain.SectorSet:
		msg.Type = domain.SectorResult
	default:
		return nil
	}
	slog.Info("INBOUND", "message", m)

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
func (f *FakeBridge) SendWithResponseChecksAndRetries(msg domain.Msg, _ func(msg domain.Msg) bool, _ time.Duration, _ int) error {
	return f.Send(msg)
}
func NewFakeBridge() *FakeBridge {
	bridge := FakeBridge{}
	slog.Info("Operating using fake bridge")
	return &bridge
}

func (f *FakeBridge) Handle() {
	return
}

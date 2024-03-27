package bridge

import (
	"context"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"log/slog"
	"time"
)

// The SerialBridge is responsible for translating commands towards things the railway can understand
type FakeBridge struct {
	broker *Broker[domain.Msg]
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

	f.broker.Send(msg)

	return nil
}
func (f *FakeBridge) SendWithResponseChecksAndRetries(msg domain.Msg, _ func(msg domain.Msg) bool, _ time.Duration, _ int) error {
	return f.Send(msg)
}
func NewFakeBridge() *FakeBridge {
	bridge := FakeBridge{
		broker: NewBroker[domain.Msg](),
	}
	slog.Info("Operating using fake bridge")
	return &bridge
}
func (f *FakeBridge) AddListener(ctx context.Context) *chan domain.Msg {
	return f.broker.AddListener(ctx)
}

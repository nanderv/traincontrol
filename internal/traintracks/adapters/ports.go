package adapters

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
)

type bridgeSender[T any] interface {
	Send(m T) error
}

type Bridge interface {
	AddReceiver(bridge.MessageReceiver)
	Send(domain.Msg) error
}

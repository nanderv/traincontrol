package core

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
)

type SendCommand interface {
	Send(m bridge.Msg) error
}

type CommandBridge interface {
	SendCommand
}

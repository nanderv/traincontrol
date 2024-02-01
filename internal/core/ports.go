package core

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
)

type SendCommand interface {
	Send(m domain.Msg) error
}

type CommandBridge interface {
	SendCommand
}

package core

import (
	"github.com/nanderv/traincontrol-prototype/internal/types"
)

type SendCommand interface {
	Send(m types.Msg)
}
type CommandBridge interface {
	SendCommand
}

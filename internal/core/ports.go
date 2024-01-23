package core

import (
	"context"
	"github.com/nanderv/traincontrol-prototype/internal/types"
)

type SendCommand interface {
	Send(m types.Msg)
}
type ReceiveCommand interface {
	BlockedReceive(ctx context.Context) (types.Msg, error)
}
type CommandBridge interface {
	SendCommand
}

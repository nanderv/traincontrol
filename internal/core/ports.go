package core

import "github.com/nanderv/traincontrol-prototype/internal/bridge/domain"

type MessageSender interface {
	Send(m domain.Msg) error
}

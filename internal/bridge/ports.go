package bridge

import "github.com/nanderv/traincontrol-prototype/internal/bridge/domain"

type MessageReceiver interface {
	Receive(msg domain.Msg) error
}

type Bridge interface {
	AddReceiver(MessageReceiver)
	Send(domain.Msg) error
}

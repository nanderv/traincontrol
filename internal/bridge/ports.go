package bridge

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"time"
)

type MessageReceiver interface {
	Receive(msg domain.Msg) error
}

type Bridge interface {
	AddReceiver(MessageReceiver)
	Send(domain.Msg) error
	SendMessageWithConfirmationAndRetries(msg domain.Msg, checker func(msg domain.Msg) bool, timeout time.Duration, retries int) error
}

package serialbridge

import (
	"github.com/nanderv/traincontrol-prototype/internal/serialbridge/domain"
	"time"
)

type MessageReceiver interface {
	Receive(msg domain.Msg) error
}

type Bridge interface {
	AddReceiver(MessageReceiver)
	Send(domain.Msg) error
	SendWithResponseChecksAndRetries(msg domain.Msg, checker func(msg domain.Msg) bool, timeout time.Duration, retries int) error
}

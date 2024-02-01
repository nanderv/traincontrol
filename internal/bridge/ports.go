package bridge

import "github.com/nanderv/traincontrol-prototype/internal/bridge/domain"

type Receiver interface {
	Receive(domain.Msg) error
}

package hwconfig

import "github.com/nanderv/traincontrol-prototype/internal/bridge/domain"

type BridgeSender[T domain.Msg] interface {
	Send(m domain.Msg) error
}

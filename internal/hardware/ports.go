package hardware

import (
	"github.com/nanderv/traincontrol-prototype/internal/hardware/domain"
)

type Sender interface {
	SetSwitchDirection(*domain.TrackSwitch, bool) error
}

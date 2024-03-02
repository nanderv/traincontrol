package traintracks

import (
	"github.com/nanderv/traincontrol-prototype/internal/traintracks/domain"
)

type Sender interface {
	SetSwitchDirection(*domain.TrackSwitch, bool) error
}

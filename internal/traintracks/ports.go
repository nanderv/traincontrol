package traintracks

type Sender interface {
	SetSwitchDirection(switchID byte, direction bool) error
}

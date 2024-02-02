package traintracks

type LayoutSender interface {
	SetSwitchDirection(switchID byte, direction bool) error
}

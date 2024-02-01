package layout

type Layout struct {
	TrackSwitches []TrackSwitch
}
type TrackSwitch struct {
	Number    byte
	Direction bool
}

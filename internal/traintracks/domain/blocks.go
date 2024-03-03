package domain

type Line struct {
	StartX int
	StartY int
	EndX   int
	EndY   int
}
type Segment struct {
	Line
	Enabled bool // Are the switches in such a way that this segment is actually running?
}
type Block struct {
	Name    string
	Segment []Segment
	Enabled bool // Is this block turned on?
}

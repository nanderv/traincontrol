package test

import (
	hwDomain "github.com/nanderv/traincontrol-prototype/internal/hardware/domain"
	layoutDomain "github.com/nanderv/traincontrol-prototype/internal/layout/domain"
	"github.com/nanderv/traincontrol-prototype/internal/serialbridge/domain"
)

func GetBaseLayout() (hwDomain.HardwareState, layoutDomain.Layout) {
	hwState := hwDomain.NewHardwareState()
	lay := layoutDomain.Layout{}
	hwState.WithTrackSwitch(hwDomain.TrackSwitch{
		Mac:       domain.Mac{22, 229, 217},
		PortID:    0,
		LeftPin:   2,
		RightPin:  3,
		Name:      "Forest",
		Direction: false,
	})
	hwState.WithTrackSwitch(hwDomain.TrackSwitch{
		Mac:       domain.Mac{22, 229, 217},
		PortID:    0,
		LeftPin:   4,
		RightPin:  5,
		Name:      "Shadow Left",
		Direction: false,
	})
	hwState.WithTrackSwitch(hwDomain.TrackSwitch{
		Mac:       domain.Mac{22, 229, 217},
		PortID:    0,
		LeftPin:   9,
		RightPin:  10,
		Name:      "Shadow Right",
		Direction: false,
	})
	hwState.WithTrackSwitch(hwDomain.TrackSwitch{
		Mac:       domain.Mac{22, 229, 217},
		PortID:    0,
		LeftPin:   11,
		RightPin:  12,
		Name:      "Loco",
		Direction: false,
	})
	lay.Blocks = append(lay.Blocks, layoutDomain.Block{
		Name: "forest_to_shadow",
		Segments: []layoutDomain.Segment{
			{
				Line: layoutDomain.Line{
					StartX: 170,
					StartY: 400,
					EndX:   200,
					EndY:   400,
				},
				Enabled: true,
			},
			{
				Line: layoutDomain.Line{
					StartX: 170,
					StartY: 380,
					EndX:   200,
					EndY:   400,
				},
				Enabled: false,
			},
			{
				Line: layoutDomain.Line{
					StartX: 200,
					StartY: 400,
					EndX:   400,
					EndY:   400,
				},
				Enabled: true,
			},
		},

		Enabled: false,
	})

	return hwState, lay
}

package main

import (
	"context"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	traintracks2 "github.com/nanderv/traincontrol-prototype/internal/bridge/adapters"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks"
	domain2 "github.com/nanderv/traincontrol-prototype/internal/traintracks/domain"
	"github.com/nanderv/traincontrol-prototype/internal/web"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	lay := domain2.NewLayout()
	lay.WithTrackSwitch(domain2.TrackSwitch{
		Mac:       domain.Mac{20, 141, 142},
		PortID:    0,
		LeftPin:   2,
		RightPin:  3,
		Name:      "1",
		Direction: false,
		X:         200,
		Y:         400,
	})
	lay.WithTrackSwitch(domain2.TrackSwitch{
		Mac:       domain.Mac{20, 141, 142},
		PortID:    0,
		LeftPin:   13,
		RightPin:  13,
		Name:      "b",
		Direction: false,
		X:         250,
		Y:         400,
	})
	lay.WithTrackSwitch(domain2.TrackSwitch{
		Mac:       domain.Mac{20, 140, 204},
		PortID:    0,
		LeftPin:   13,
		RightPin:  13,
		Name:      "sl",
		Direction: false,
		X:         450,
		Y:         400,
	})
	lay.WithTrackSwitch(domain2.TrackSwitch{
		Mac:       domain.Mac{22, 229, 217},
		PortID:    0,
		LeftPin:   13,
		RightPin:  13,
		Name:      "sl2",
		Direction: false,
		X:         450,
		Y:         350,
	})
	lay.WithBlock(domain2.Block{
		Name: "forest_to_shadow",
		Segments: []domain2.Segment{
			{
				Line: domain2.Line{
					StartX: 170,
					StartY: 400,
					EndX:   200,
					EndY:   400,
				},
				Enabled: true,
			},
			{
				Line: domain2.Line{
					StartX: 170,
					StartY: 380,
					EndX:   200,
					EndY:   400,
				},
				Enabled: false,
			},
			{
				Line: domain2.Line{
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

	lay.WithBlock(domain2.Block{
		Name: "forest_siding",
		Segments: []domain2.Segment{

			{
				Line: domain2.Line{
					StartX: 110,
					StartY: 340,
					EndX:   167,
					EndY:   378,
				},
				Enabled: true,
			},
		},

		Enabled: false,
	})
	lay.WithBlock(domain2.Block{
		Name: "forest_curve",
		Segments: []domain2.Segment{

			{
				Line: domain2.Line{
					StartX: 110,
					StartY: 400,
					EndX:   167,
					EndY:   400,
				},
				Enabled: true,
			},
		},

		Enabled: false,
	})

	lay.WithBlock(domain2.Block{
		Name: "forest_to_main",
		Segments: []domain2.Segment{

			{
				Line: domain2.Line{
					StartX: 50,
					StartY: 400,
					EndX:   107,
					EndY:   400,
				},
				Enabled: true,
			},
		},

		Enabled: false,
	})
	lay.WithBlock(domain2.Block{
		Name: "shadow_left",
		Segments: []domain2.Segment{

			{
				Line: domain2.Line{
					StartX: 410,
					StartY: 400,
					EndX:   450,
					EndY:   400,
				},
				Enabled: true,
			},
			{
				Line: domain2.Line{
					StartX: 450,
					StartY: 400,
					EndX:   480,
					EndY:   400,
				},
				Enabled: true,
			},
			{
				Line: domain2.Line{
					StartX: 450,
					StartY: 400,
					EndX:   480,
					EndY:   370,
				},
				Enabled: true,
			},
		},

		Enabled: false,
	})
	lay.WithBlock(domain2.Block{
		Name: "shadow_top_l",
		Segments: []domain2.Segment{

			{
				Line: domain2.Line{
					StartX: 483,
					StartY: 370,
					EndX:   557,
					EndY:   370,
				},
				Enabled: true,
			},
		},
		Enabled: false,
	})
	lay.WithBlock(domain2.Block{
		Name: "shadow_bottom_l",
		Segments: []domain2.Segment{

			{
				Line: domain2.Line{
					StartX: 483,
					StartY: 400,
					EndX:   557,
					EndY:   400,
				},
				Enabled: true,
			},
		},
		Enabled: false,
	})
	lay.WithBlock(domain2.Block{
		Name: "shadow_top_r",
		Segments: []domain2.Segment{

			{
				Line: domain2.Line{
					StartX: 560,
					StartY: 370,
					EndX:   634,
					EndY:   370,
				},
				Enabled: true,
			},
		},
		Enabled: false,
	})
	lay.WithBlock(domain2.Block{
		Name: "shadow_bottom_r",
		Segments: []domain2.Segment{

			{
				Line: domain2.Line{
					StartX: 560,
					StartY: 400,
					EndX:   634,
					EndY:   400,
				},
				Enabled: true,
			},
		},
		Enabled: false,
	})
	c, err := traintracks.NewTrackService(lay)

	b := bridge.NewSerialBridge()

	traintracks2.NewMessageAdapter(c, b)

	if err != nil {
		return
	}
	go func() {
		err := web.Init(context.Background(), c)
		if err != nil {
			fmt.Println(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	x := <-done
	slog.Info("Application killed", "signal", x.String())
}

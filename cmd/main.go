package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/hardware"
	hwDomain "github.com/nanderv/traincontrol-prototype/internal/hardware/domain"
	"github.com/nanderv/traincontrol-prototype/internal/serialbridge"
	traintracks2 "github.com/nanderv/traincontrol-prototype/internal/serialbridge/adapters"
	"github.com/nanderv/traincontrol-prototype/internal/serialbridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/web"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func getBaseLayout() hwDomain.HardwareState {
	lay := hwDomain.NewHardwareState()
	lay.WithTrackSwitch(hwDomain.TrackSwitch{
		Mac:       domain.Mac{20, 141, 142},
		PortID:    0,
		LeftPin:   2,
		RightPin:  3,
		Name:      "1",
		Direction: false,
	})
	lay.WithTrackSwitch(hwDomain.TrackSwitch{
		Mac:       domain.Mac{20, 141, 142},
		PortID:    0,
		LeftPin:   13,
		RightPin:  13,
		Name:      "b",
		Direction: false,
	})
	lay.WithTrackSwitch(hwDomain.TrackSwitch{
		Mac:       domain.Mac{20, 140, 204},
		PortID:    0,
		LeftPin:   13,
		RightPin:  13,
		Name:      "sl",
		Direction: false,
	})
	lay.WithTrackSwitch(hwDomain.TrackSwitch{
		Mac:       domain.Mac{22, 229, 217},
		PortID:    0,
		LeftPin:   13,
		RightPin:  13,
		Name:      "sl2",
		Direction: false,
	})
	lay.WithBlock(hwDomain.Block{
		Name: "forest_to_shadow",
		Segments: []hwDomain.Segment{
			{
				Line: hwDomain.Line{
					StartX: 170,
					StartY: 400,
					EndX:   200,
					EndY:   400,
				},
				Enabled: true,
			},
			{
				Line: hwDomain.Line{
					StartX: 170,
					StartY: 380,
					EndX:   200,
					EndY:   400,
				},
				Enabled: false,
			},
			{
				Line: hwDomain.Line{
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

	lay.WithBlock(hwDomain.Block{
		Name: "forest_siding",
		Segments: []hwDomain.Segment{

			{
				Line: hwDomain.Line{
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
	lay.WithBlock(hwDomain.Block{
		Name: "forest_curve",
		Segments: []hwDomain.Segment{

			{
				Line: hwDomain.Line{
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

	lay.WithBlock(hwDomain.Block{
		Name: "forest_to_main",
		Segments: []hwDomain.Segment{

			{
				Line: hwDomain.Line{
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
	lay.WithBlock(hwDomain.Block{
		Name: "shadow_left",
		Segments: []hwDomain.Segment{

			{
				Line: hwDomain.Line{
					StartX: 410,
					StartY: 400,
					EndX:   450,
					EndY:   400,
				},
				Enabled: true,
			},
			{
				Line: hwDomain.Line{
					StartX: 450,
					StartY: 400,
					EndX:   480,
					EndY:   400,
				},
				Enabled: true,
			},
			{
				Line: hwDomain.Line{
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
	lay.WithBlock(hwDomain.Block{
		Name: "shadow_top_l",
		Segments: []hwDomain.Segment{

			{
				Line: hwDomain.Line{
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
	lay.WithBlock(hwDomain.Block{
		Name: "shadow_bottom_l",
		Segments: []hwDomain.Segment{

			{
				Line: hwDomain.Line{
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
	lay.WithBlock(hwDomain.Block{
		Name: "shadow_top_r",
		Segments: []hwDomain.Segment{

			{
				Line: hwDomain.Line{
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
	lay.WithBlock(hwDomain.Block{
		Name: "shadow_bottom_r",
		Segments: []hwDomain.Segment{

			{
				Line: hwDomain.Line{
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

	return lay
}
func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	brdg := "fake"
	flag.StringVar(&brdg, "bridge", "fake", "Set which hardware brdige to use (fake or serial)")
	if brdg != "fake" && brdg != "serial" {
		panic("Bridge ain't real")
	}

	c, err := hardware.NewTrackService(getBaseLayout())

	var b serialbridge.Bridge
	if brdg == "fake" {
		b = serialbridge.NewFakeBridge()
	}
	if brdg == "serial" {
		b = serialbridge.NewSerialBridge()
	}
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

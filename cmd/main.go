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
	})

	lay.WithBlock(domain2.Block{
		Name: "testBlock",
		Segment: []domain2.Segment{
			{
				Line: domain2.Line{
					StartX: 0,
					StartY: 0,
					EndX:   0,
					EndY:   0,
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

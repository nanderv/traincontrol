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
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	switches := map[string]*domain2.TrackSwitch{
		"1": {
			Mac:       domain.Mac{20, 141, 142},
			PortID:    0,
			LeftPin:   2,
			RightPin:  3,
			Name:      "1",
			Direction: false,
		},
	}

	c, err := traintracks.NewTrackService(domain2.Layout{TrackSwitches: switches})

	bridg := bridge.NewSerialBridge()
	go bridg.IncomingHandler()
	go bridg.OutgoingHandler()

	traintracks2.NewMessageAdapter(c, bridg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err != nil {
		return
	}
	go func() {
		err := web.Init(ctx, c)
		if err != nil {
			fmt.Println(err)
		}
	}()
	for {
		time.Sleep(5 * time.Second)
		c.SetSwitchDirection("1", true)
		time.Sleep(5 * time.Second)
		c.SetSwitchDirection("1", false)
	}
	time.Sleep(1 * time.Hour)
}

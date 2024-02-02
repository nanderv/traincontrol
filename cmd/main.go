package main

import (
	"context"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks/adapters"
	"github.com/nanderv/traincontrol-prototype/internal/web"
	"log/slog"
	"os"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	c, err := traintracks.NewCore(traintracks.WithTrackSwitch(1), traintracks.WithTrackSwitch(2), traintracks.WithTrackSwitch(3))

	bridg := bridge.NewSerialBridge()
	go bridg.Handle()

	adapters.NewMessageAdapter(c, bridg)

	//hwConf := hwconfig.HwConfigurator{}
	//hwAdapters.NewMessageAdapter(&hwConf, bridg)

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

	err = c.SetSwitchAction(1, true)
	if err != nil {
		slog.Error("Could not set switch", "error", err)
		//return
	}
	err = c.SetSwitchAction(1, true)
	if err != nil {
		slog.Error("Could not set switch", "error", err)
		//return
	}

	time.Sleep(1 * time.Hour)
}

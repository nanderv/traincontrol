package main

import (
	"context"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	hwAdapters "github.com/nanderv/traincontrol-prototype/internal/bridge/adapters/hwconfig"
	traintracks2 "github.com/nanderv/traincontrol-prototype/internal/bridge/adapters/traintracks"
	"github.com/nanderv/traincontrol-prototype/internal/hwconfig"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks"
	"github.com/nanderv/traincontrol-prototype/internal/web"
	"log/slog"
	"os"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	c, err := traintracks.NewCore(traintracks.WithTrackSwitch(1), traintracks.WithTrackSwitch(2), traintracks.WithTrackSwitch(3))

	bridg := bridge.NewSerialBridge()
	go bridg.IncomingHandler()
	go bridg.OutgoingHandler()

	traintracks2.NewMessageAdapter(c, bridg)

	hwConf := hwconfig.NewHWConfigurator()
	hwAdapters.NewMessageAdapter(hwConf, bridg)

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

	time.Sleep(1 * time.Hour)
}

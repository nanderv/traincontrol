package main

import (
	"context"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/core"
	"github.com/nanderv/traincontrol-prototype/internal/core/adapters"
	"github.com/nanderv/traincontrol-prototype/internal/web"
	"log/slog"
	"os"
	"time"
)

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	bridg := bridge.NewSerialBridge()

	go bridg.Handle()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := core.NewCore(core.WithTrackSwitch(1), core.WithTrackSwitch(2), core.WithTrackSwitch(3))

	adapters.NewMessageAdapter(c, bridg)
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

package main

import (
	"context"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/core"
	"github.com/nanderv/traincontrol-prototype/internal/web"
	"log/slog"
	"os"
	"time"
)

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := core.NewCore(core.WithBridge(), core.WithTrackSwitch(1), core.WithTrackSwitch(2), core.WithTrackSwitch(3))
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

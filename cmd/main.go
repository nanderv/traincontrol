package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/datasets/test"
	"github.com/nanderv/traincontrol-prototype/internal/hardware"
	"github.com/nanderv/traincontrol-prototype/internal/serialbridge"
	traintracks2 "github.com/nanderv/traincontrol-prototype/internal/serialbridge/adapters"
	"github.com/nanderv/traincontrol-prototype/internal/web"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	brdg := "fake"
	flag.StringVar(&brdg, "bridge", "fake", "Set which hardware brdige to use (fake or serial)")
	if brdg != "fake" && brdg != "serial" {
		panic("Bridge ain't real")
	}
	layStr := "test"
	flag.StringVar(&layStr, "layout", "test", "Set which layout to use")
	c, err := hardware.NewTrackService(test.GetBaseLayout())

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

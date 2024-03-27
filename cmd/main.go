package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/datasets/test"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/hardware"
	bridgeAdapter "github.com/nanderv/traincontrol-prototype/internal/hardware/adapters"
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

	layStr := "test"
	flag.StringVar(&layStr, "layout", "test", "Set which layout to use")
	flag.Parse()

	if brdg != "fake" && brdg != "serial" {
		panic("Bridge ain't real")
	}

	hw, _ := test.GetBaseLayout()
	c, err := hardware.NewTrackService(hw)

	var b bridgeAdapter.Bridge
	fmt.Println(brdg)
	if brdg == "fake" {
		b = bridge.NewFakeBridge()
	}
	if brdg == "serial" {
		fmt.Println("HI")
		b = bridge.NewSerialBridge()
	}
	bridgeAdapter.NewAdapter(c, b)

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

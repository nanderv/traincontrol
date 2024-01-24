package main

import (
	"context"
	"github.com/nanderv/traincontrol-prototype/internal/core"
	"github.com/nanderv/traincontrol-prototype/internal/http_adapter"
	"time"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	newCore, err := core.NewCore(core.WithFakeBridge(), core.WithTrackSwitch(1), core.WithTrackSwitch(2), core.WithTrackSwitch(3))
	if err != nil {
		return
	}
	go http_adapter.Init()

	err = newCore.SetSwitchAction(1, true)
	if err != nil {
		panic(err)
	}
	newCore.SetSwitchAction(2, true)
	newCore.SetSwitchAction(2, false)

	time.Sleep(1 * time.Hour)

	cancel()

}

package main

import (
	"github.com/nanderv/traincontrol-prototype/internal/core"
	"github.com/nanderv/traincontrol-prototype/internal/http_adapter"
	"time"
)

func main() {
	core, err := core.NewCore(core.WithFakeBridge(), core.WithTrackSwitch(1), core.WithTrackSwitch(2), core.WithTrackSwitch(3))
	if err != nil {
		return
	}
	go http_adapter.Init()

	err = core.SetSwitchAction(1, true)
	if err != nil {
		panic(err)
	}
	core.SetSwitchAction(2, true)
	core.SetSwitchAction(2, false)

	time.Sleep(1 * time.Hour)
}

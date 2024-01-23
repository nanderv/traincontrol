package main

import (
	"context"
	"github.com/nanderv/traincontrol-prototype/internal/core"
	"github.com/nanderv/traincontrol-prototype/internal/http_adapter"
	"time"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	newCore, err := core.NewCore(core.WithFakeBridge())
	if err != nil {
		return
	}
	go http_adapter.Init()

	newCore.SetSwitch(1, true)
	newCore.SetSwitch(2, true)
	newCore.SetSwitch(2, false)

	time.Sleep(1 * time.Hour)

	cancel()

}

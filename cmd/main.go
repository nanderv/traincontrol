package main

import (
	"context"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/core"
	"github.com/nanderv/traincontrol-prototype/internal/web"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := core.NewCore(core.WithFakeBridge(), core.WithTrackSwitch(1), core.WithTrackSwitch(2), core.WithTrackSwitch(3))
	if err != nil {
		return
	}
	go func() {
		err := web.Init(ctx, c)
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(5 * time.Second)
	err = c.SetSwitchAction(1, true)
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)

	c.SetSwitchAction(2, true)
	time.Sleep(5 * time.Second)

	c.SetSwitchAction(2, false)

	time.Sleep(1 * time.Hour)
}

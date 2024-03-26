package main

import (
	"flag"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/datasets/test"
	"github.com/nanderv/traincontrol-prototype/internal/hardware"
	"github.com/nanderv/traincontrol-prototype/internal/serialbridge"
	traintracks2 "github.com/nanderv/traincontrol-prototype/internal/serialbridge/adapters"
	"github.com/rivo/tview"
	"log/slog"
)

type NulWriter struct{}

func (w *NulWriter) Write([]byte) (int, error) { return 0, nil }
func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(&NulWriter{}, nil)))
	brdg := "fake"
	flag.StringVar(&brdg, "bridge", "fake", "Set which hardware brdige to use (fake or serial)")
	if brdg != "fake" && brdg != "serial" {
		panic("Bridge ain't real")
	}
	layStr := "test"
	flag.StringVar(&layStr, "layout", "test", "Set which layout to use")
	lay := test.GetBaseLayout()
	c, err := hardware.NewTrackService(lay)
	if err != nil {
		panic(err)
	}

	var b serialbridge.Bridge
	if brdg == "fake" {
		b = serialbridge.NewFakeBridge()
	}
	if brdg == "serial" {
		b = serialbridge.NewSerialBridge()
	}
	traintracks2.NewMessageAdapter(c, b)
	cha := c.AddNewReturnChannel()
	app := tview.NewApplication()
	form := tview.NewForm()
	buttons := make(map[string]*tview.Button)
	for _, sw := range lay.TrackSwitches {
		s := sw
		but := tview.NewButton(fmt.Sprintf("ST %s %v", s.Name, s.Direction))
		fmt.Println(but)
		buttons[s.Name] = but
	}
	go func() {
		select {
		case sta := <-*cha:
			lay = sta
			i := 0
			for _, sw := range sta.TrackSwitches {
				s := sw

				form.GetButton(i).SetLabel(fmt.Sprintf("ST %s %v", s.Name, s.Direction))

				i++
			}
		}
	}()
	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

package main

import (
	"flag"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/datasets/test"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	traintracks2 "github.com/nanderv/traincontrol-prototype/internal/bridge/adapters"
	"github.com/nanderv/traincontrol-prototype/internal/hardware"
	"github.com/rivo/tview"
	"log/slog"
)

type NulWriter struct {
	Res *chan string
}

func (w *NulWriter) Write(b []byte) (int, error) {
	bb := string(b)
	go func() { *w.Res <- bb }()
	return len(bb), nil
}
func main() {
	nw := NulWriter{}
	ch := make(chan string)
	nw.Res = &ch
	slog.SetDefault(slog.New(slog.NewJSONHandler(&nw, nil)))

	brdg := "fake"
	flag.StringVar(&brdg, "bridge", "fake", "Set which hardware brdige to use (fake or serial)")

	layStr := "test"
	flag.StringVar(&layStr, "layout", "test", "Set which layout to use")
	flag.Parse()

	if brdg != "fake" && brdg != "serial" {
		panic("Bridge ain't real")
	}
	hwState, _ := test.GetBaseLayout()
	c, err := hardware.NewTrackService(hwState)
	if err != nil {
		panic(err)
	}

	var b bridge.Bridge
	if brdg == "fake" {
		b = bridge.NewFakeBridge()
	}
	if brdg == "serial" {
		b = bridge.NewSerialBridge()
	}
	traintracks2.NewMessageAdapter(c, b)

	cha := c.AddNewReturnChannel()
	app := tview.NewApplication()
	grid := tview.NewFlex()
	form := tview.NewForm()
	form.SetTitle("Switch Control")
	grid = grid.AddItem(form, 5, 30, true)
	txt := tview.NewTextView().SetDynamicColors(true).
		SetRegions(true)
	txt.SetTitle("Logs")

	txt.SetBorder(true)
	grid = grid.AddItem(txt, 0, 70, true).SetDirection(tview.FlexRow)
	buttons := make(map[string]*tview.Button)
	for _, sw := range hwState.TrackSwitches {
		ss := sw
		form.AddButton(ss.Name, func() {
			s, err := hwState.GetSwitch(ss.Name)
			if err != nil {
				panic(err)
			}
			go c.SetSwitchDirection(s.Name, !s.Direction)
		})
		but := form.GetButton(form.GetButtonIndex(sw.Name))
		form.SetButtonsAlign(1)

		buttons[sw.Name] = but
	}

	go func() {
		for {

			select {

			case sta := <-*cha:
				app.QueueUpdate(func() {

					for _, sw := range sta.TrackSwitches {
						if sw.Direction {
							buttons[sw.Name].SetLabel(fmt.Sprintf("S %s  %v", sw.Name, sw.Direction))
						} else {
							buttons[sw.Name].SetLabel(fmt.Sprintf("S %s %v", sw.Name, sw.Direction))
						}

					}
					app.ForceDraw()
				})

			case w := <-*nw.Res:
				app.QueueUpdate(func() {
					txt.SetText(txt.GetText(true) + w)
					app.ForceDraw()
				})
			}
		}
	}()
	form.SetBorder(true).SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

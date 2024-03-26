package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nanderv/traincontrol-prototype/internal/hardware"
	"github.com/nanderv/traincontrol-prototype/internal/hardware/domain"
	"io"
	"log/slog"
)

type LayoutHTTPAdapter struct {
	c  *hardware.TrackService
	ch *chan domain.HardwareState
	h  io.Writer
}

func NewLayoutHTTP(c *hardware.TrackService, h io.Writer) *LayoutHTTPAdapter {
	return &LayoutHTTPAdapter{
		c:  c,
		ch: c.AddNewReturnChannel(),
		h:  h,
	}
}

func (l *LayoutHTTPAdapter) Handle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case d := <-*l.ch:
			html := ""
			i := 0
			for _, y := range d.TrackSwitches {
				var dir = "-"
				if y.Direction {
					dir = "x"
				}
				i += 1
				html +=
					fmt.Sprintf(`<div style="position:absolute; left: %vpx; top: %vpx; border: 3px solid #73AD21;" onClick="btn('%s', %v)">%s</div>`, 50*i, 100, y.Name, !y.Direction, dir)
			}
			//panic("disco")
			//html = "<div>" + html + "</div>"
			//if err != nil {
			//	slog.Error("unable to marshall data!", "error", err)
			//}

			_, err := l.h.Write([]byte(html))
			if err != nil {
				slog.Error("Unable to send out data to writer!", "error", err)
			}
		}
	}
}

func (l *LayoutHTTPAdapter) SetSwitch(c echo.Context) error {
	input := setSwitchInput{}

	err := json.NewDecoder(c.Request().Body).Decode(&input)
	if err != nil {
		return err
	}

	return l.c.SetSwitchDirection(input.SwitchID, input.Direction)
}

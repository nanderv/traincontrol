package web

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks/domain"
	"io"
	"log/slog"
)

type LayoutJSONAdapter struct {
	c  *traintracks.TrackService
	ch *chan domain.Layout
	h  io.Writer
}

func NewLayoutJSONAdapter(c *traintracks.TrackService, h io.Writer) *LayoutJSONAdapter {
	return &LayoutJSONAdapter{
		c:  c,
		ch: c.AddNewReturnChannel(),
		h:  h,
	}
}

func (l *LayoutJSONAdapter) Handle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case d := <-*l.ch:
			b, err := json.Marshal(&d)
			if err != nil {
				slog.Error("unable to marshall data!", "error", err)
			}

			_, err = l.h.Write(b)
			if err != nil {
				slog.Error("Unable to send out data to writer!", "error", err)
			}
		}
	}
}

type setSwitchInput struct {
	SwitchID  string `json:"switch_id"`
	Direction bool   `json:"direction"`
}

func (l *LayoutJSONAdapter) SetSwitch(c echo.Context) error {
	input := setSwitchInput{}

	err := json.NewDecoder(c.Request().Body).Decode(&input)
	if err != nil {
		return err
	}

	return l.c.SetSwitchDirection(input.SwitchID, input.Direction)
}

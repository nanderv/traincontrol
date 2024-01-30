package web

import (
	"context"
	"encoding/json"
	"github.com/nanderv/traincontrol-prototype/internal/core"
	"io"
)

type LayoutAdapter struct {
	c  *core.Core
	ch *chan core.State
	h  io.Writer
}

func NewLayoutAdapter(c *core.Core, h io.Writer) *LayoutAdapter {
	return &LayoutAdapter{
		c:  c,
		ch: c.AddNewReturnChannel(),
		h:  h,
	}
}

func (l *LayoutAdapter) Handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():

			return nil
		case d := <-*l.ch:
			b, err := json.Marshal(&d)
			if err != nil {
				return err
			}

			_, err = l.h.Write(b)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

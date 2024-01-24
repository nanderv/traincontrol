package http_adapter

import (
	"context"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/core"
)

type LayoutAdapter struct {
	c  *core.Core
	ch *chan struct{}
	h  *MessageRouter[routeMessage]
}

func NewLayoutAdapter(c *core.Core, ch *chan struct{}, h *MessageRouter[routeMessage]) *LayoutAdapter {
	return &LayoutAdapter{
		c:  c,
		ch: ch,
		h:  h,
	}
}

func (l *LayoutAdapter) Handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-*l.ch:
			fmt.Println("HI")
		}
	}

}

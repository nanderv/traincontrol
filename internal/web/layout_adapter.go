package web

import (
	"context"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/core"
	"io"
	"time"
)

type LayoutAdapter struct {
	c  *core.Core
	ch *chan struct{}
	h  io.Writer
}

func NewLayoutAdapter(c *core.Core, ch *chan struct{}, h io.Writer) *LayoutAdapter {
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
			_, err := l.h.Write([]byte("hi"))
			if err != nil {
				return err
			}
			fmt.Println("HI")
		}
		time.Sleep(100 * time.Millisecond)
	}

}

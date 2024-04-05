package bridge

import (
	"context"
	"log/slog"
	"sync"
)

type Broker[T any] struct {
	sync.Mutex
	clients map[*client[T]]struct{}
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{

		clients: make(map[*client[T]]struct{}),
	}
}

func (b *Broker[T]) AddListener(ctx context.Context) *chan T {
	b.Lock()
	defer b.Unlock()
	ic := make(chan T)
	oc := make(chan T)
	ls := make([]T, 0)
	c := client[T]{inputChannel: &ic, outputChannel: &oc, buffer: ls}
	go c.handle(ctx, b)
	b.clients[&c] = struct{}{}
	return c.outputChannel
}

func (b *Broker[T]) Send(msg T) {
	b.Lock()
	defer b.Unlock()
	for client := range b.clients {
		*client.inputChannel <- msg
	}
}

type client[T any] struct {
	inputChannel  *chan T
	buffer        []T
	outputChannel *chan T
}

func (c *client[T]) handle(ctx context.Context, b *Broker[T]) {
	c.handleLoop(ctx)
	b.Lock()
	defer b.Unlock()
	slog.Debug("Cleanup")
	delete(b.clients, c)
	// explicit cleanup
	close(*c.inputChannel)
	close(*c.outputChannel)
}

func (c *client[T]) handleLoop(ctx context.Context) {
	for {
		if len(c.buffer) > 0 {
			select {
			case <-ctx.Done():
				return
			case m := <-*c.inputChannel:
				c.buffer = append(c.buffer, m)

			case *c.outputChannel <- c.buffer[0]:
				c.buffer = c.buffer[1:]
			}
		} else {
			select {
			case <-ctx.Done():
				return
			case m := <-*c.inputChannel:
				c.buffer = append(c.buffer, m)
			}
		}
	}
}

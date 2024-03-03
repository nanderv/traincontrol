package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"sync"
)

type MessageRouter struct {
	sync.RWMutex
	channelMap map[*chan []byte]struct{}
}

func (r *MessageRouter) Subscribe() *chan []byte {
	r.Lock()
	defer r.Unlock()
	slog.Debug("Subscribing to channel")
	c := make(chan []byte)
	r.channelMap[&c] = struct{}{}
	return &c
}

func (r *MessageRouter) Unsubscribe(c *chan []byte) {
	r.Lock()
	defer r.Unlock()
	slog.Debug("Unsubscribing from channel")

	delete(r.channelMap, c)
}

func (r *MessageRouter) Write(in []byte) (int, error) {
	r.RLock()
	defer r.RUnlock()

	for c, _ := range r.channelMap {
		*c <- in
	}
	return len(in), nil
}

func NewRouter() *MessageRouter {
	m := make(map[*chan []byte]struct{})
	return &MessageRouter{
		channelMap: m,
	}
}

func RouteWithMessageRouter(router *MessageRouter) echo.HandlerFunc {
	return func(c echo.Context) error {
		cha := router.Subscribe()
		defer router.Unsubscribe(cha)

		w := c.Response().Writer
		flusher, ok := w.(http.Flusher)

		if !ok {
			panic("expected web.ResponseWriter to be an web.Flusher")
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Transfer-Encoding", "chunked")
		for {
			select {
			case t := <-*cha:
				_, _ = fmt.Fprintln(w, "event: update")
				_, _ = fmt.Fprintf(w, "data: %s\n\n", t)

				flusher.Flush()
			case <-c.Request().Context().Done():
				return nil
			}
		}
	}
}

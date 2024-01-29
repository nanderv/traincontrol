package web

import (
	"fmt"
	"net/http"
)

type MessageRouter[T fmt.Stringer] struct {
	channelMap map[*chan T]struct{}
}

func (r *MessageRouter[T]) Subscribe() *chan T {
	fmt.Println("Connect")
	c := make(chan T)
	r.channelMap[&c] = struct{}{}
	return &c
}
func (r *MessageRouter[T]) Unsubscribe(c *chan T) {
	fmt.Println("Disco")
	delete(r.channelMap, c)
}
func (r *MessageRouter[T]) Send(in T) {
	for c, _ := range r.channelMap {
		*c <- in
	}
}
func NewRouter[T fmt.Stringer]() *MessageRouter[T] {
	m := make(map[*chan T]struct{})
	return &MessageRouter[T]{
		channelMap: m,
	}
}

func RouteWithMessageRouter[T fmt.Stringer](router *MessageRouter[T]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c := router.Subscribe()
		defer router.Unsubscribe(c)

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
			case t := <-*c:
				fmt.Println(t)

				_, err := fmt.Fprintln(w, "event: update")
				if err != nil {
					return
				}

				_, err = fmt.Fprintf(w, "data: %s\n\n", t)

				if err != nil {
					return
				}
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	}
}

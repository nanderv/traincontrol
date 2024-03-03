package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks"
	"log"
	"net/http"
)

type in struct {
	Method   string `json:"data"`
	Envelope string `json:"envelope"`
}

type routeMessage string

func (m routeMessage) String() string {
	return string(m)
}

func act(router *MessageRouter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Println("Wrong method")
			_, err := w.Write([]byte("Wrong method"))
			if err != nil {
				return
			}
		}
		var v in
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			fmt.Println("EE", err)
			return
		}
	}
}

func Init(ctx context.Context, c *traintracks.TrackService) error {
	// Add file server
	fs := http.FileServer(http.Dir("internal/web/static/"))
	http.Handle("/", http.StripPrefix("/", fs))
	// Add route for getting chunked data
	rt := NewRouter()
	go NewLayoutAdapter(c, rt).Handle(ctx)

	http.HandleFunc("/send", act(rt))
	http.HandleFunc("/chunk", RouteWithMessageRouter(rt))

	// Start the server
	log.Print("Listening on localhost:9898")
	log.Fatal(http.ListenAndServe(":9898", nil))
	return nil
}

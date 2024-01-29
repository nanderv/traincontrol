package web

import (
	"encoding/json"
	"fmt"
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

func act(router *MessageRouter[routeMessage]) func(w http.ResponseWriter, r *http.Request) {
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
func Init() {
	// Add file server
	fs := http.FileServer(http.Dir("webroot/"))
	http.Handle("/", http.StripPrefix("/", fs))

	// Add route for getting chunked data
	rt := NewRouter[routeMessage]()
	http.HandleFunc("/send", act(rt))
	http.HandleFunc("/chunk", RouteWithMessageRouter(rt))

	// Start the server
	log.Print("Listening on localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

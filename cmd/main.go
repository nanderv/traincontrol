package main

import "github.com/nanderv/traincontrol-prototype/internal/http_adapter"

func main() {
	go http_adapter.Init()
	select {}
}

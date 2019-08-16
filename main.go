package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/codeka/mcsvradmin/config"
	"github.com/codeka/mcsvradmin/mcsvr"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello  df " + message
	w.Write([]byte(message))
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	inst, err := mcsvr.Launch(cfg)
	if err != nil {
		log.Fatalf("error launching: %v", err)
	}
	inst.MonitorProcess()
	// Stuff.

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ping", sayHello)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

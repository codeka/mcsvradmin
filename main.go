package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/codeka/mcsvradmin/config"
	"github.com/codeka/mcsvradmin/mcsvr"
	"github.com/codeka/mcsvradmin/static"
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

	http.Handle("/", http.FileServer(static.Assets))
        http.HandleFunc("/ping", sayHello)

        url := fmt.Sprintf(":%d", cfg.ListenPort)
        log.Printf("Web server listening on '%s'...", url)
        if err := http.ListenAndServe(url, nil); err != nil {
                panic(err)
        }
}


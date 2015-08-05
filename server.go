package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kju2/buildbulb/job"
	"github.com/kju2/buildbulb/light"
	"github.com/kju2/buildbulb/notification"
	"github.com/kju2/buildbulb/util"
)

// TODO As a server administrator I want to configure the server with cmd arguments.

func main() {
	port := "8080"
	notifier, notifications := notification.NewController()
	_, status := job.NewController(notifications)
	_, err := light.NewController(status)

	if err != nil {
		log.Fatal(err)
	}

	util.Log.WithField("port", port).Info("Will listen forever for HTTP requests.")
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/notify", notifier.Handle)
	http.ListenAndServe(":"+port, nil)
}

func pong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "pong")
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/kju2/buildbulb/job"
	"github.com/kju2/buildbulb/light"
	"github.com/kju2/buildbulb/notification"
	"github.com/kju2/buildbulb/util"
)

var (
	port         = flag.Int("port", 8080, "port to listen for incoming HTTP requests")
	jobsFilePath = flag.String("jobsFilePath", "~/.buildbulb", "path to load and persist jobs")
)

func main() {
	flag.Parse()

	notifier, notifications := notification.NewController()
	_, status := job.NewController(notifications, *jobsFilePath)
	_, err := light.NewController(status)

	if err != nil {
		log.Fatal(err)
	}

	util.Log.WithField("port", port).Info("Will listen forever for HTTP requests.")
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/notify", notifier.Handle)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}

func pong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "pong")
}

package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kju2/buildbulb/job"
	"github.com/kju2/buildbulb/light"
	"github.com/kju2/buildbulb/notification"
	"github.com/kju2/buildbulb/util"
)

var (
	port         = flag.Int("port", 8080, "port to listen for incoming HTTP requests")
	jobsFilePath = flag.String("jobsFilePath", "", "path to load and persist jobs")
)

func main() {
	flag.Parse()

	if len(*jobsFilePath) == 0 {
		util.Log.Warn("Path to load and persist jobs hasn't been provided.")
	}

	notifier, notifications := notification.NewController()
	jobifier, status := job.NewController(notifications, *jobsFilePath)
	_, err := light.NewController(status)

	if err != nil {
		util.Log.WithField("error", err).Fatal("Light controller threw an error on initialization.")
	}

	util.Log.WithField("port", *port).Info("Will listen forever for HTTP requests.")
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/jobs", jobifier.Handle)
	http.HandleFunc("/notify", notifier.Handle)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}

func pong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "pong")
}

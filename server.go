package main

import (
	"log"
	"net/http"

	"github.com/kju2/buildbulb/light"
	"github.com/kju2/buildbulb/notification"
	"github.com/kju2/buildbulb/project"
)

func main() {
	notifier, notifications := notification.NewController()
	_, status := project.NewController(notifications)
	_, err := light.NewController(status)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/notify", notifier.Handle)
	http.ListenAndServe(":8080", nil)
}

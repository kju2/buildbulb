package notification

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"

	"github.com/kju2/buildbulb/job"
	"github.com/kju2/buildbulb/util"
)

type Controller struct {
	output chan<- *job.Job
}

func NewController() (*Controller, <-chan *job.Job) {
	output := make(chan *job.Job, 1)
	return &Controller{output}, output
}

func (c *Controller) Handle(w http.ResponseWriter, r *http.Request) {
	requestDump, _ := httputil.DumpRequest(r, true)
	util.Log.Debugf("Received request: %s", requestDump)

	var msg notification
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		badRequest(w, "Error occured parsing request: '"+err.Error()+"'.")
		return
	}

	job, err := msg.job()
	if err != nil {
		badRequest(w, "Error occured parsing request: '"+err.Error()+"'.")
		return
	}

	util.Log.WithField("job", *job).Info("Decoded job.")
	c.output <- job

	success(w)
}

func badRequest(w http.ResponseWriter, error string) {
	util.Log.Info(error)
	http.Error(w, error, http.StatusBadRequest)
}

func success(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

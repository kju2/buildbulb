package notification

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"

	"github.com/kju2/buildbulb/project"
	"github.com/kju2/buildbulb/util"
)

type Controller struct {
	output chan<- *project.Project
}

func NewController() (*Controller, <-chan *project.Project) {
	output := make(chan *project.Project, 1)
	return &Controller{output}, output
}

func (c *Controller) Handle(w http.ResponseWriter, r *http.Request) {
	requestDump, _ := httputil.DumpRequest(r, true)
	util.Log.Debugf("Received request: %s", requestDump)

	var job job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		badRequest(w, "Error occured parsing request: '"+err.Error()+"'.")
		return
	}

	project, err := job.project()
	if err != nil {
		badRequest(w, "Error occured parsing request: '"+err.Error()+"'.")
		return
	}

	util.Log.Infof(`Decoded job: "%+v"`, *project)
	c.output <- project

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

package notification

import (
	"fmt"
	"time"

	"github.com/kju2/buildbulb/job"
)

type notification struct {
	Name  string
	Build build
}

type build struct {
	Status string
}

func (msg *notification) job() (*job.Job, error) {
	if len(msg.Name) < 1 {
		return nil, fmt.Errorf("Job name is missing.")
	}

	status, err := job.Parse(msg.Build.Status)
	if err != nil {
		return nil, err
	}

	return job.NewJob(msg.Name, status, time.Now().Round(time.Second)), nil
}

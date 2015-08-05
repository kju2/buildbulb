package job

import (
	"fmt"
	"strings"
)

type Job struct {
	Name   string
	Status Status
}

func NewJob(name string, status Status) *Job {
	return &Job{name, status}
}

// The build status of a job.
type Status int

const (
	Failure  Status = iota // Job couldn't be compiled or another unrecoverable error occurred.
	Unstable               // At least one test for this job failed.
	Success                // Job compiled and all tests are green.
)

func (s Status) String() string {
	switch s {
	case Failure:
		return "Failure"
	case Unstable:
		return "Unstable"
	case Success:
		return "Success"
	default:
		return "Unknown"
	}
}

func Parse(status string) (Status, error) {
	if len(status) < 1 {
		return Failure, fmt.Errorf("Given status is empty.")
	}

	got, ok := map[string]Status{"failure": Failure, "unstable": Unstable, "success": Success}[strings.ToLower(status)]
	if !ok {
		return Failure, fmt.Errorf("Given status %q is invalid.", status)
	}
	return got, nil
}

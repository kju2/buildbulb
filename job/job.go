package job

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Job struct {
	Name        string
	Status      Status
	LastUpdated time.Time
}

func NewJob(name string, status Status, lastUpdated time.Time) *Job {
	return &Job{name, status, lastUpdated}
}

// The build status of a job.
type Status int

const (
	Failure  Status = iota // Job couldn't be compiled or another unrecoverable error occurred.
	Unstable               // At least one test for this job failed.
	Success                // Job compiled and all tests are green.
)

func (job Job) String() string {
	return fmt.Sprintf("[%s:%s]", job.Name, job.Status);
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Status) UnmarshalJSON(b []byte) error {
	var err error = nil
	*s, err = Parse(string(b))
	return err
}

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

	// If the source of the received status is the JSON document, it contains quotation marks which have to be removed.
	parsedStatus := strings.ToLower(strings.Replace(status, "\"", "", 2))
	
	got, ok := map[string]Status{"failure": Failure, "unstable": Unstable, "success": Success}[parsedStatus]
	if !ok {
		return Failure, fmt.Errorf("Given status %q is invalid.", status)
	}

	return got, nil
}

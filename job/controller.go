package job

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/kju2/buildbulb/util"
)

type Controller struct {
	jobs         []*Job
	jobsFilePath string
	mutex        sync.RWMutex
	output       chan<- Status
}

func NewController(input <-chan *Job, jobsFilePath string) (*Controller, <-chan Status) {
	output := make(chan Status)
	jobs := make([]*Job, 0)
	c := &Controller{jobs, jobsFilePath, sync.RWMutex{}, output}

	go c.run(input)

	return c, output
}

func (c *Controller) run(input <-chan *Job) {
	c.readJobsFromFile()
	c.sendOverallStatus()

	for job := range input {
		c.updateJob(job)
		c.writeJobsToFile()
		c.sendOverallStatus()
	}
}

func (c *Controller) updateJob(job *Job) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	defer util.Log.WithField("jobs", c.jobs).Debug("State of all jobs.")

	for i, j := range c.jobs {
		if j.Name == job.Name {
			c.jobs[i] = job
			return
		}
	}

	c.jobs = append(c.jobs, job)
}

func (c *Controller) overallStatus() Status {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Determine overall job status.
	overallStatus := Success
	for _, job := range c.jobs {
		switch job.Status {
		case Failure:
			util.Log.WithField("job", job).Error("Job failed.")
		case Unstable:
			util.Log.WithField("job", job).Warn("Job is unstable.")
		}

		if job.Status < overallStatus {
			overallStatus = job.Status
		}
	}
	util.Log.WithField("overall status", overallStatus).Info("Overall status determined.")

	return overallStatus
}

func (c *Controller) sendOverallStatus() {
	c.output <- c.overallStatus()
}

func (c *Controller) writeJobs(w io.Writer) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	b, err := json.MarshalIndent(c.jobs, "", "\t")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = json.Indent(&buf, b, "", "\t")
	if err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}

func (c *Controller) writeJobsToFile() {
	if len(c.jobsFilePath) == 0 {
		return
	}

	f, err := os.Create(c.jobsFilePath)
	if err != nil {
		util.Log.WithField("error", err).Error("Couldn't open file to persist jobs.")
		return
	}
	defer f.Close()

	err = c.writeJobs(f)
	if err != nil {
		util.Log.WithField("error", err).Error("Couldn't persist jobs.")
		return
	}
}

func (c *Controller) readJobsFromFile() {
	if len(c.jobsFilePath) == 0 {
		return
	}

	f, err := os.Open(c.jobsFilePath)
	if err != nil {
		util.Log.WithField("error", err).Warn("Couldn't open file with persisted jobs.")
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	if err := json.NewDecoder(f).Decode(&c.jobs); err != nil {
		util.Log.WithField("error", err).Error("Couldn't parse persisted jobs.")
		return
	}

	util.Log.WithField("jobs", c.jobs).Info("Read persisted Jobs.")
}

func (c *Controller) Handle(w http.ResponseWriter, r *http.Request) {
	err := c.writeJobs(w)
	if err != nil {
		util.Log.WithField("error", err).Error("Error occurred encoding jobs.")
		http.Error(w, "Couldn't marshal jobs", http.StatusInternalServerError)
	}
}

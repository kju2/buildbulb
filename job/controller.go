package job

import (
	"encoding/json"
	"github.com/kju2/buildbulb/util"
	"os"
)

type Controller struct {
	jobs         map[string]Status
	jobsFilePath string
	output       chan<- Status
}

func NewController(input <-chan *Job, jobsFilePath string) (*Controller, <-chan Status) {
	output := make(chan Status)
	jobs := make(map[string]Status)
	c := &Controller{jobs, jobsFilePath, output}

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
	c.jobs[job.Name] = job.Status
	util.Log.WithField("jobs", c.jobs).Debug("State of all jobs.")
}

func (c *Controller) sendOverallStatus() {
	c.output <- c.overallStatus()
}

func (c *Controller) overallStatus() Status {
	// Determine overall job status.
	overallStatus := Success
	for job, status := range c.jobs {
		switch status {
		case Failure:
			util.Log.WithField("job", job).Error("Job failed.")
		case Unstable:
			util.Log.WithField("job", job).Warn("Job is unstable.")
		}

		if status < overallStatus {
			overallStatus = status
		}
	}
	util.Log.WithField("overall status", overallStatus).Info("Overall status determined.")

	return overallStatus
}

func (c *Controller) writeJobsToFile() {
	b, err := json.MarshalIndent(c.jobs, "", "\t")
	if err != nil {
		util.Log.WithField("error", err).Error("Couldn't persist jobs.")
		return
	}

	f, err := os.Create(c.jobsFilePath)
	if err != nil {
		util.Log.WithField("error", err).Error("Couldn't persist jobs.")
		return
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		util.Log.WithField("error", err).Error("Couldn't persist jobs.")
		return
	}
}

func (c *Controller) readJobsFromFile() {
	f, err := os.Open(c.jobsFilePath)
	if err != nil {
		util.Log.WithField("error", err).Warn("Couldn't find persisted jobs.")
		return
	}

	if err := json.NewDecoder(f).Decode(&c.jobs); err != nil {
		util.Log.WithField("error", err).Error("Couldn't parse persisted jobs.")
		return
	}

	util.Log.WithField("jobs", c.jobs).Info("Read persisted Jobs.")
}

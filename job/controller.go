package job

import "github.com/kju2/buildbulb/util"

type Controller struct {
	jobs   map[string]Status
	output chan<- Status
}

func NewController(input <-chan *Job) (*Controller, <-chan Status) {
	output := make(chan Status)
	jobs := make(map[string]Status)
	c := &Controller{jobs, output}

	go c.run(input)

	return c, output
}

func (c *Controller) run(input <-chan *Job) {
	for job := range input {
		// Update status of job.
		c.jobs[job.Name] = job.Status
		util.Log.WithField("jobs", c.jobs).Debug("State of all jobs.")

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

		// Send overall status to the receiver.
		c.output <- overallStatus

		// Persist state of jobs.
		// TODO
	}
}

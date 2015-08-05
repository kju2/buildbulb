package project

import (
	"fmt"
	"strings"

	"github.com/kju2/buildbulb/util"
)

type Project struct {
	Name   Name
	Status Status
}

func NewProject(name Name, status Status) *Project {
	return &Project{name, status}
}

// The name of a project.
type Name string

// The build status of a project.
type Status int

const (
	Failure  Status = iota // Project couldn't be compiled or another unrecoverable error occurred.
	Unstable               // At least one test for this project failed.
	Success                // Project compiled and all tests are green.
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

type Controller struct {
	projects map[Name]Status
	output   chan<- Status
}

func NewController(input <-chan *Project) (*Controller, <-chan Status) {
	output := make(chan Status)
	projects := make(map[Name]Status)
	c := &Controller{projects, output}

	go c.run(input)

	return c, output
}

func (c *Controller) run(input <-chan *Project) {
	for notification := range input {
		// Update project status.
		newStatus := notification.Status
		oldStatus, _ := c.projects[notification.Name]
		util.Log.Infof("Updating project %q from %q to %q", notification.Name, oldStatus, newStatus)

		c.projects[notification.Name] = notification.Status
		util.Log.Debugf("Projects: %+v", c.projects)

		// Determine overall project status.
		overallStatus := Success
		for project, status := range c.projects {
			switch status {
			case Failure:
				util.Log.Errorf("Project %q failed.", project)
			case Unstable:
				util.Log.Warnf("Project %q is unstable.", project)
			}

			if status < overallStatus {
				overallStatus = status
			}
		}

		// Send overall status to the receiver.
		c.output <- overallStatus

		// Persist state of projects.
		// TODO
	}
}

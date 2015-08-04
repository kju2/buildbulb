package project

import (
	"fmt"
	"log"
	"strings"
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
	Failure Status = iota
	Unstable
	Success
)

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
		log.Printf("Notification: %+v\n", notification)

		c.projects[notification.Name] = notification.Status
		log.Printf("Projects: %+v\n", c.projects)

		overallStatus := Success
		for _, status := range c.projects {
			if status < overallStatus {
				overallStatus = status
			}
		}
		c.output <- overallStatus
	}
}

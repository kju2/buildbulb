package notification

import (
	"fmt"

	"github.com/kju2/buildbulb/project"
)

type job struct {
	Name  string
	Build build
}

type build struct {
	Status string
}

func (j *job) project() (*project.Project, error) {
	if len(j.Name) < 1 {
		return nil, fmt.Errorf("Project name is missing.")
	}
	name := project.Name(j.Name)

	status, err := project.Parse(j.Build.Status)
	if err != nil {
		return nil, err
	}

	return project.NewProject(name, status), nil
}

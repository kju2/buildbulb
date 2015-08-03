package notification

import (
	"fmt"
	"strings"

	"github.com/kju2/buildbulb/project"
)

type job struct {
	Name  string
	Build build
}

func (j *job) isFinalized() bool {
	_, found := phases[strings.ToLower(j.Build.Phase)]
	return found
}

func (j *job) project() (*project.Project, error) {
	if len(j.Name) < 1 {
		return nil, fmt.Errorf("Project name is missing in '%+v'", j)
	}
	name := project.Name(j.Name)

	status, err := project.Parse(j.Build.Status)
	if err != nil {
		return nil, err
	}

	return project.NewProject(name, status), nil
}

//func (m *message) UnmarshalJSON(data []byte) error {
//var aux struct {
//Name  string
//Build build
//}
//if err := json.NewDecoder(bytes.NewReader(data)).Decode(&aux); err != nil {
//return err
//}
//log.Printf("Unmarshal Message: %+v\n", aux)
//m.Name = aux.Name
//m.Build = aux.Build
//return nil
//}

type build struct {
	Phase  string
	Status string
}

var (
	phases = map[string]bool{"completed": true, "finished": true, "finalized": true}
)

//func (b *build) UnmarshalJSON(data []byte) error {

//}

//func (b *build) xUnmarshalJSON(data []byte) error {
//var build map[string]interface{}

//if err := json.Unmarshal(data, &build); err != nil {
//return err
//}

//rawStatus, found := build["status"]
//if found == false {
//return fmt.Errorf("status not contained in notification")
//}
//status, err := project.UnmarshalJSON(rawStatus.(string))
//if err != nil {
//return err
//}
//n.Status = status
//return nil
//}

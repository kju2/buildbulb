package notification

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kju2/buildbulb/project"
)

const notification_template = `{
    "name": "%s",
    "url": "job/#{project}/",
    "build": {
        "phase": "COMPLETED",
        "status": "%s"
    }
}`

func TestDecodingJob(t *testing.T) {
	projectName := "test_project"
	status := "Success"

	p, err := decodeJob(createTestInput(projectName, status))
	if err != nil {
		t.Fatal(err)
	}

	if string(p.Name) != projectName || p.Status != project.Success {
		t.Fatal("Decoded project doesn't have the right values")
	}
}

// As a user of this package I want to get an error if a build doesn't have a status.
func TestDecodingJobWithoutStatus(t *testing.T) {
	input := strings.NewReader(`{"name": "test", "build": {}}`)
	_, err := decodeJob(input)
	if err == nil {
		t.Fatal("Status is missing, but job could still be parsed.")
	}
}

// As a user of this package I want to get an error if a build has an empty status.
func TestDecodingJobWithEmptyStatus(t *testing.T) {
	input := strings.NewReader(`{"name": "test", "build": {"status": ""}}`)
	_, err := decodeJob(input)
	if err == nil {
		t.Fatal("Status is empty, but job could still be decoded.")
	}
}

// As a user of this package I want to get an error if a build has an invalid status.
func TestDecodingJobWithAnInvalidStatus(t *testing.T) {
	input := strings.NewReader(`{"name": "test", "build": {"status": "blub"}}`)
	_, err := decodeJob(input)
	if err == nil {
		t.Fatal("Status is empty, but job could still be decoded.")
	}
}
func TestHttpHandling(t *testing.T) {
	c, output := NewController()

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "", createTestInput("test_project", "Failure"))
	if err != nil {
		t.Error(err)
	}

	c.Handle(response, request)
	<-output
}

func createTestInput(name, status string) io.Reader {
	return strings.NewReader(fmt.Sprintf(notification_template, name, status))
}

func decodeJob(r io.Reader) (*project.Project, error) {
	var j job
	if err := json.NewDecoder(r).Decode(&j); err != nil {
		return nil, err
	}
	return j.project()
}

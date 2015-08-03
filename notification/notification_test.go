package notification

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const notification_template = `{
    "name": "%s",
    "url": "job/#{project}/",
    "build": {
        "full_url": "http://localhost:8080/job/#{project}/48/",
        "number": 48,
        "phase": "FINALIZED",
        "status": "%s",
        "url": "job/#{project}/48/",
        "scm": {
            "url": "git@github.com:jenkinsci/#{project}.git",
            "branch": "origin/master",
            "commit": "4886d1ff4821879410f4f4a93168e6cc179a8eb3"
        },
        "artifacts": {
            "test.jar": [
                "http://localhost:8080/job/test/48/artifact/target/test.jar"
            ],
            "test.hpi": [
                "http://localhost:8080/job/test/48/artifact/target/test.hpi"
            ]
        }
    }
}`

func TestDecodingNotification(t *testing.T) {
	projectName := "test_project"
	status := "Success"

	n, err := decodeMessage(createTestInput(projectName, status))
	if err != nil {
		t.Error(err)
	}
	log.Printf("HERE: %+v\n", n)
	//if n.Name != projectName || n.Status != project.Success {
	//t.FailNow()
	//}
}

func TestHttpHandling(t *testing.T) {
	c, output := NewController()

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "", createTestInput("BlaProject", "Success"))
	if err != nil {
		t.Error(err)
	}

	c.Handle(response, request)
	//var n *Notification
	var p = <-output
	//log.Printf("Response: %+v\n", response.Code)
	log.Printf("HERE: %+v\n", p)
}

func createTestInput(name, status string) io.Reader {
	return strings.NewReader(fmt.Sprintf(notification_template, name, status))
}

func decodeMessage(r io.Reader) (*job, error) {
	var j job
	if err := json.NewDecoder(r).Decode(&j); err != nil {
		return nil, err
	}
	return &j, nil
}

func TestIfJobIsFinished(t *testing.T) {
	isFinalized("finished", t)
	isFinalized("Finished", t)
	isFinalized("FINISHED", t)
}

func TestIfJobIsCompleted(t *testing.T) {
	isFinalized("completed", t)
	isFinalized("Completed", t)
	isFinalized("COMPLETED", t)
}

func TestIfJobIsFinalized(t *testing.T) {
	isFinalized("finalized", t)
	isFinalized("Finalized", t)
	isFinalized("FINALIZED", t)
}

func TestIfJobIsNotFinalized(t *testing.T) {
	isNotFinalized("", t)
	isNotFinalized("bla", t)
	isNotFinalized("test Completed", t)
	isNotFinalized(" Completed", t)
	isNotFinalized("Completed test", t)
}
func isFinalized(phase string, t *testing.T) {
	j := &job{}
	j.Build.Phase = phase

	if !j.isFinalized() {
		t.Fatalf("Job phase '%v' should be finalized, but isn't", phase)
	}
}

func isNotFinalized(phase string, t *testing.T) {
	j := &job{}
	j.Build.Phase = phase

	if j.isFinalized() {
		t.Fatalf("Job phase '%v' should be finalized, but isn't", phase)
	}
}

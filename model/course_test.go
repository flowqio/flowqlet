package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestModelCourse(t *testing.T) {

	course := &Course{}

	data, _ := json.Marshal(course)

	t.Log(string(data))
}

func TestModelScenario(t *testing.T) {

	scenario := &Scenario{}

	scenario.Title = "Docker start"
	scenario.Comments = "This is about docker start"
	scenario.URL = "docker-start-overview"
	scenario.Backend = Backend{Environment: []string{"Test=123", "PWD_IP=10.0.0.1"}, Image: "stevensu/dind", Expose: []string{"80/tcp", "8080/tcp"}, Mounts: []string{`type=bind,source=/home/opc/data/{uid},target=/app,readonly`, `source=/home/opc/data123,target=/app1`}}
	scenario.Created = Time(time.Now())
	scenario.Intro = MarkDown{Title: "intro", File: "intro.md"}

	scenario.Steps = []MarkDown{
		MarkDown{Title: "This is step01", File: "step01.md"},
		MarkDown{Title: "This is step02", File: "step02.md"},
	}
	ParsePortBidings(scenario.Backend)
	t.Log(scenario.Backend.GetMounts(map[string]string{"uid": "1000"}))
	data, err := json.MarshalIndent(scenario, "", " ")

	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

//ComposeUP used create and run container
func ComposeUP(oid, scenario string, parameter ...map[string]string) ([]*Instance, error) {

	_parametr := map[string]string{"OID": oid}

	if len(parameter) > 0 {
		for k, v := range parameter[0] {
			_parametr[k] = v
		}

	}

	projectName := scenario
	if len(scenario) > 22 {
		projectName = scenario[:22]
	}

	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles:      []string{configPath(scenario)},
			ProjectName:       projectName,
			TemplateParameter: _parametr,
		},
	}, nil)

	if err != nil {
		return nil, err
	}

	err = project.Up(context.Background(), options.Up{})

	if err != nil {
		return nil, err
	}

	containers := project.GetContainerInfo()

	_instances := []*Instance{}
	now := time.Now()
	duration, _ := time.ParseDuration(instanceMaxLive) //instanceMaxLive最大生命周期

	for idx := range containers {
		instance := &Instance{ID: containers[idx].ID, Name: containers[idx].Name, CreatedAt: now, ExpiresAt: now.Add(duration), OID: oid, ScenarioName: scenario, ExposePorts: exposePorts(containers[idx].Ports)}
		_instances = append(_instances, instance)
		if _, ok := instances[oid]; !ok {
			instances[oid] = make(map[string]*Instance)
		}
		instances[oid][containers[idx].ID] = instance

	}

	cli := GetEtcdClient()

	data, err := json.Marshal(_instances)

	if err != nil {
		return nil, err
	}

	cli.Put(context.TODO(), fmt.Sprintf("/flowq/user/%s/%s", oid, scenario), string(data))

	return _instances, nil

}

//ComposeDown stop and remove container from docker-compose.yml
func ComposeDown(oid, scenario string, parameter ...map[string]string) error {

	_parametr := map[string]string{"OID": oid}

	if len(parameter) > 0 {
		for k, v := range parameter[0] {
			_parametr[k] = v
		}

	}
	projectName := scenario
	if len(scenario) > 22 {
		projectName = scenario[:22]
	}

	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles:      []string{configPath(scenario)},
			ProjectName:       projectName,
			TemplateParameter: _parametr,
		},
	}, nil)

	if err != nil {
		return err
	}

	err = project.Down(context.Background(), options.Down{})

	if err != nil {
		return err
	}
	for k, _ := range instances {
		if k == oid {
			for k1, v1 := range instances[k] {
				if v1.ScenarioName == scenario {
					delete(instances[k], k1)
				}
			}

			if len(instances[k]) == 0 {
				delete(instances, k)
			}
		}
	}

	cli := GetEtcdClient()

	cli.Delete(context.TODO(), fmt.Sprintf("/flowq/user/%s/%s", oid, scenario))

	return nil

}

func configPath(scenario string) string {
	return fmt.Sprintf("scenario/%s/docker-compose.yml", scenario)
}

func exposePorts(raw string) []string {
	if raw == "" {
		return []string{}
	}

	ports := strings.Split(raw, ",")
	result := []string{}

	for idx := range ports {
		fmt.Println(ports[idx])
		if len(strings.Split(ports[idx], "->")) > 1 {
			result = append(result, strings.TrimSpace(strings.Split(ports[idx], "->")[1]))
		}

	}

	return result

}

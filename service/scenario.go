package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func ClearScenario(uid string) error {

	cli := GetEtcdClient()
	_, err := cli.Delete(context.TODO(), fmt.Sprintf("/flowq/user/%s", uid), clientv3.WithPrefix())
	return err

}

func UpdateScenarioFile(scenario string) error {
	timeFormat := "200601021504"
	ts := time.Now().Format(timeFormat)
	scenarioPath := configPath(scenario)

	if UpdateServerEndpoint() == "" {

		return fmt.Errorf("need setup update server ")
	}
	//back scenario
	err := os.Rename(scenarioPath, scenarioPath+"."+ts)
	if err != nil {
		return err
	}

	out, err := os.Create(scenarioPath)
	if err != nil {
		return err
	}
	defer out.Close()
	resource := fmt.Sprintf(UpdateServerEndpoint()+"/%s/docker-compose.yml", scenario)
	resp, err := http.Get(resource)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil

}

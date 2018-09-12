// Copyright 2018 flowq Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"

	"crypto/tls"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

var c *client.Client

const (
	Byte     = 1
	Kilobyte = 1024 * Byte
	Megabyte = 1024 * Kilobyte
)

type UInt16Slice []uint16

type Instance struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	Hostname     string                  `json:"hostname,omitempty"`
	IP           string                  `json:"ip,omitempty"`
	conn         *types.HijackedResponse `json:"-"`
	ctx          context.Context         `json:"-"`
	dockerClient *client.Client          `json:"-"`
	IsManager    *bool                   `json:"is_manager,omitempty"`
	Mem          string                  `json:"mem,omitempty"`
	Cpu          string                  `json:"cpu,omitempty"`
	Ports        UInt16Slice             `json:"-"`
	tempPorts    []uint16                `json:"-"`
	ServerCert   []byte                  `json:"server_cert,omitempty"`
	ServerKey    []byte                  `json:"server_key,omitempty"`
	cert         *tls.Certificate        `json:"-"`
	CreatedAt    time.Time               `json:"created_at,omitempty"`
	ExpiresAt    time.Time               `json:"expires_at,omitempty"`
	Image        string                  `json:"image,omitempty"`
	ExposePorts  []string                `json:"expose,omitempty"`
	ScenarioName string                  `json:"scenario,omitempty"`
	OID          string                  `json:"oid,omitempty"`
}

func init() {
	var err error
	c, err = client.NewEnvClient()
	if err != nil {
		// this wont happen if daemon is offline, only for some critical errors
		log.Fatal("Cannot initialize docker client")
	}

	gob.Register(Instance{})
}

func CreateAttachConnection(id string, ctx context.Context) (*types.HijackedResponse, error) {

	conf := types.ContainerAttachOptions{true, true, true, true, "ctrl-^,ctrl-^", true}
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	conn, err := c.ContainerAttach(ctx, id, conf)
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

func CheckContainer(id string) (types.ContainerJSON, bool) {
	c, err := client.NewEnvClient()
	detail, err := c.ContainerInspect(context.Background(), id)
	if err != nil {
		return detail, false
	}
	return detail, true
}

func CreateExecAttachConnection(ctx context.Context, id, cmd string) (*types.HijackedResponse, error) {

	if cmd == "" {
		cmd = "sh"
	}
	execConfig := types.ExecConfig{
		User:         "root",
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		Cmd:          strings.Split(cmd, " "),
		Tty:          true,
		Detach:       false,
		DetachKeys:   "ctrl-^,ctrl-^",
	}

	c, err := client.NewEnvClient()
	if err != nil {
		log.Println("1", err)
		return nil, err
	}
	log.Println(id, execConfig)

	execInstance, err := c.ContainerExecCreate(ctx, id, execConfig)
	log.Println(execInstance, err)
	log.Printf("ContainerExec %s, %s ,%v", id, cmd, execInstance)
	conn, err := c.ContainerExecAttach(ctx, execInstance.ID, execConfig)
	if err != nil {
		log.Println("container attach", err)
		return nil, err
	}

	return &conn, nil

}

// func NewInstance(uid string, backend model.Backend) (*Instance, error) {

// 	if backend.Image == "" {
// 		backend.Image = "stevensu/dind"
// 	}

// 	return CreateInstance(backend, uid)
// }

func CloseInstance(uid, id string) error {
	var durationMilliseconds time.Duration = 500 * time.Millisecond

	err := c.ContainerStop(context.Background(), id, &durationMilliseconds)
	if err != nil {
		if strings.Index(fmt.Sprintf("%v", err), "No such container") == -1 {
			return err
		}

	}

	delete(instances[uid], id)

	if len(instances[uid]) == 0 {
		delete(instances, uid)
		ClearNetwork(id)
	}

	//Save Instance to Disk
	SaveInstanceToDisk()

	log.Println("Clear instance success")
	return nil
}

func GetInstance(uid ...string) interface{} {
	if len(uid) == 0 {
		return instances
	}
	return instances[uid[0]]
}

func GetContainerInfo(id string) (types.ContainerJSON, error) {
	return c.ContainerInspect(context.Background(), id)
}

func ClearNetwork(id string) error {

	err := c.NetworkRemove(context.Background(), id)

	if err != nil {
		log.Printf("Clear network %s err [%s]\n", id, err)
		return err
	}
	return nil
}

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
	"errors"
	"os"

	"github.com/coreos/etcd/clientv3"
	log "github.com/sirupsen/logrus"

	"github.com/flowqio/flowqlet/config"
)

var serviceRoot = "/flowq/env"

var tokenRoot = "/flowq/token"

var apiToken = ""

var REQUEST_TIMEOUT = 5

var nodeConfig *config.FlowqletConfig

//GetFlowLetConfig return config.FlowqletConfig
func GetFlowqLetConfig() config.FlowqletConfig {
	return *nodeConfig
}

//OnBoard flowqlet will be put nodeinfo to etcdserver
func OnBoard(token string, nodeID string, letInfo string) error {

	cli := GetEtcdClient()

	k, err := cli.Get(context.TODO(), tokenRoot)
	if err != nil || len(k.Kvs) == 0 {
		log.Fatalf("Token [%s] not found,join failure", token)
		os.Exit(-1)
	}

	if string(k.Kvs[0].Value) != token {
		log.Fatalf("Token [%s] not math,join failure", token)
		os.Exit(-1)
	}

	apiToken = token

	key := serviceRoot + "/" + nodeID

	k, err = cli.Get(context.TODO(), key)

	if err != nil || len(k.Kvs) > 0 {
		log.Fatalf("NodeID [%s] already online   %s", key, k.Kvs)
	}

	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = cli.Put(context.TODO(), key, string(letInfo), clientv3.WithLease(resp.ID))

	log.Infof("Lease ID: %+v", resp.ID)

	if err != nil {
		log.Fatal(err)
		return err
	}

	ch, err := cli.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for {
		select {
		case <-cli.Ctx().Done():
			log.Fatal(errors.New("server closed"))
		case _, ok := <-ch:
			if !ok {
				log.Error("keep alive channel closed")
				OnBoard(token, nodeID, letInfo)
				return nil
			} else {
				log.Debugf("Recv reply from service: %s, ttl:%d", nodeID, ka.TTL)
			}
		}
	}

}

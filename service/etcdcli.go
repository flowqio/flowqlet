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
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	log "github.com/sirupsen/logrus"
)

var _etcdv3Client *clientv3.Client
var _etcdEndporint = "127.0.0.1:2379"

func InitEtcdClient(endporint ...string) (err error) {

	if len(endporint) > 0 {
		_etcdEndporint = endporint[0]
	}
	_etcdv3Client, err = clientv3.New(clientv3.Config{
		DialTimeout: 2 * time.Second,
		Endpoints:   strings.Split(_etcdEndporint, ","),
		Username:    "steven",
		Password:    "welcome1",
	})
	if err == nil {
		log.Debugf("Init etcd client sucess")
	}

	return err
}

func GetEtcdClient() *clientv3.Client {
	return _etcdv3Client
}

package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/flowqio/flowqlet/config"
	log "github.com/sirupsen/logrus"
)

func init() {

	log.Debug("Clear Scheduler start success")
	LoadInstanceFromDisk()
	go func() {
		watch()
	}()
}

var updateServerEndpoint = ""

func UpdateServerEndpoint(endpoint ...string) string {
	if len(endpoint) > 0 {
		updateServerEndpoint = endpoint[0]
	}
	return updateServerEndpoint
}

func InitFlowqlet(nodeConf *config.FlowqletConfig) {

	nodeConfig = nodeConf

	InitEtcdClient(nodeConf.EtcdURL)

	UpdateServerEndpoint(nodeConf.UpdateServer)

	//go service.OnBoard(*token, *nodeID, fmt.Sprintf("%s:%d", *host, *port))

	go OnBoard(nodeConf.Token, nodeConf.NodeID, nodeConf.Addr)

	srv := &http.Server{
		Addr: nodeConf.Addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("flowqlet-%s listen at %s ", nodeConf.NodeID, nodeConf.Addr)
	if nodeConf.UpdateServer != "" {
		log.Printf("UpdateServer: %s", nodeConf.UpdateServer)
	}
	if len(nodeConf.Label) > 0 {
		log.Printf("Label: %s", strings.Join(nodeConf.Label, ","))
	}
	log.Fatal(srv.ListenAndServe())
}

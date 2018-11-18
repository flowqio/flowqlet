package config

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

type FlowqletConfig struct {
	NodeID       string   `json:"nodeID,omitempty"`
	Addr         string   `json:"addr,omitempty"`
	Token        string   `json:"-"`
	Label        []string `json:"label,omitempty"`
	EtcdURL      string   `json:"etcdURL,omitempty"`
	UpdateServer string   `json:"updateServer,omitempty"`
	VolumeDriver string   `json:"voludmeDriver,omitempty"`
}

var nodeConfig *FlowqletConfig

func NodeConfig() FlowqletConfig {
	return *nodeConfig
}
func InitFlag() (*FlowqletConfig, error) {

	host := flag.String("host", "localhost", "listen address")
	port := flag.Int("port", 8800, "listen port")
	nodeID := flag.String("nodeID", "env1", "nodeID")
	token := flag.String("token", "", "join token")
	etcdEndpoint := flag.String("etcdURL", "127.0.0.1:2379", "etcd endporint")
	updateserver := flag.String("updateServer", "", "update web server endporint")
	volumeDriver := flag.String("volumeDriver", "local", "volume driver support local|s3")
	label := flag.String("label", "", "label")

	flag.Parse()

	var err error

	if *token == "" {
		return nil, fmt.Errorf("join token required.")
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", *host, *port))

	if err != nil {
		return nil, err
	}

	nodeConfig, err = NewFlowqletConfig(*nodeID, tcpAddr.String(), *token, *etcdEndpoint, *updateserver, *volumeDriver, *label)
	if err != nil {
		return nil, err
	}

	return nodeConfig, nil
}

//NewFlowletConfig return FlowqletConfig
func NewFlowqletConfig(nodeID, addr, token, etcdURL, updateServer, volumeDriver, label string) (*FlowqletConfig, error) {

	labels := []string{}
	if label != "" {
		labels = strings.Split(label, ",")
	}
	return &FlowqletConfig{
		NodeID:       nodeID,
		Addr:         addr,
		Token:        token,
		Label:        labels,
		EtcdURL:      etcdURL,
		UpdateServer: updateServer,
		VolumeDriver: volumeDriver,
	}, nil
}

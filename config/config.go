package config

import (
	"os"
	"flag"
	"fmt"
	"net"
	"strings"
	"github.com/flowqio/flowqlet/version"
)

//FlowqletConfig is config struct 
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

//NodeConfig return *nodeConfig FlowletConfig
func NodeConfig() FlowqletConfig {
	return *nodeConfig
}

//InitFlag init service used FlowqletConfig
func InitFlag() (*FlowqletConfig, error) {

    //all flag include host, port , nodeID,token,etcdEndpoint,updateserver,volumeDriver
	host := flag.String("host", "localhost", "listen address")
	port := flag.Int("port", 8800, "listen port")
	nodeID := flag.String("nodeID", "env1", "nodeID")
	token := flag.String("token", "", "join token")
	etcdEndpoint := flag.String("etcdURL", "127.0.0.1:2379", "etcd endporint")
	updateServer := flag.String("updateServer", "", "update web server endporint")
	volumeDriver := flag.String("volumeDriver", "local", "volume driver support local|s3")
	label := flag.String("label", "", "label")
	_version:=flag.Bool("version",false,"show version")

	flag.Parse()

	if *_version {
		version.PrintVersionInfo()
		os.Exit(-1)
	}

	var err error

	if *token == "" {
		return nil, fmt.Errorf("join token required.")
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", *host, *port))

	if err != nil {
		return nil, err
	}

	nodeConfig, err = NewFlowqletConfig(*nodeID, tcpAddr.String(), *token, *etcdEndpoint, *updateServer, *volumeDriver, *label)
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

package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	//Init router handler prefix and utils log configuration
	_ "github.com/flowqio/flowqlet/router"

	"github.com/flowqio/flowqlet/service"
)

func main() {

	port := flag.Int("port", 8800, "listen port")
	host := flag.String("host", "localhost", "listen address")
	nodeID := flag.String("n", "env1", "nodeID")
	token := flag.String("token", "", "join token")
	etcdEndpoint := flag.String("etcdurl", "127.0.0.1:2379", "etcd endporint")
	flag.Parse()

	if *token == "" {
		fmt.Println("join token required.")
		return
	}
	service.InitEtcdClient(*etcdEndpoint)

	go service.OnBoard(*token, *nodeID, fmt.Sprintf("%s:%d", *host, *port))

	srv := &http.Server{
		Addr: *host + ":" + strconv.Itoa(*port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("flowqlet [%s] listen at %s:%d ", *nodeID, *host, *port)
	log.Fatal(srv.ListenAndServe())

}

package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/flowqio/flowqlet/config"
	"github.com/flowqio/flowqlet/service"
	"github.com/flowqio/flowqlet/version"

	//init router handler configuration
	_ "github.com/flowqio/flowqlet/router"
)

func main() {

	//print version
	version.PrintBanner()

	//init flag config
	conf, err := config.InitFlag()
	if err != nil {
		log.Fatal(err)
	}

	//init service
	service.InitFlowqlet(conf)

}

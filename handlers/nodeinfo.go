package handlers

import (
	"net/http"

	"github.com/flowqio/flowqlet/config"
)

func NodeInfo(w http.ResponseWriter, r *http.Request) {

	ServerJSON(w, config.NodeConfig())

}

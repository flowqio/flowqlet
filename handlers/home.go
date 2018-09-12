package handlers

import (
	"net/http"

	"github.com/flowqio/flowqlet/version"
)

func Index(w http.ResponseWriter, r *http.Request) {

	ServerJSON(w, map[string]interface{}{"status": "ok", "version": version.VersionInfo()})
	return

}

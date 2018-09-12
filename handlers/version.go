package handlers

import (
	"net/http"

	"github.com/flowqio/flowqlet/version"
)

//APIVersion /api/v1 return information
func APIVersion(w http.ResponseWriter, r *http.Request) {
	ServerJSON(w, version.VersionInfo())
}

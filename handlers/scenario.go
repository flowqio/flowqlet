package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/flowqio/flowqlet/version"
	"github.com/gorilla/mux"

	"github.com/flowqio/flowqlet/service"
)

func UpdateScenario(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	sid, ok := vars["sid"]

	if !ok || sid == "" {
		msg, _ := json.Marshal(map[string]string{"msg": "missing SID", "status": "400"})
		http.Error(w, string(msg), http.StatusBadRequest)
		return
	}

	if !service.ScenarioIsExits(sid) {
		msg, _ := json.Marshal(map[string]string{"msg": "file not exits", "status": "400"})
		http.Error(w, string(msg), http.StatusBadRequest)
		return
	}

	err := service.UpdateScenarioFile(sid)
	if err != nil {
		msg, _ := json.Marshal(map[string]string{"msg": err.Error(), "status": "400"})
		http.Error(w, string(msg), http.StatusBadRequest)
		return
	}
	ServerJSON(w, map[string]interface{}{"status": "ok", "version": version.VersionInfo(), "sid": sid})
	return

}

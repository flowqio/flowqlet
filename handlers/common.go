package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//ServerJSON is warpper func help handler service json content
func ServerJSON(w http.ResponseWriter, model interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "flowqlet")
	w.Header().Set("Vendor", "flowq.io")

	err := json.NewEncoder(w).Encode(model)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
	}
}

func ConsumeJSON(r *http.Request, model interface{}) error {
	contentWriter := json.NewDecoder(r.Body)
	err := contentWriter.Decode(model)
	if err != nil {
		return err
	}
	return nil
}

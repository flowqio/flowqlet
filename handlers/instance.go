package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"

	"github.com/flowqio/flowqlet/service"
)

//persistent path

/* ----- old func , new version use libcompose ------ */
// func CreateInstance(w http.ResponseWriter, r *http.Request) {

// 	var oid = ""

// 	vars := mux.Vars(r)

// 	if _, ok := vars["oid"]; ok {
// 		oid = vars["oid"]
// 	} else {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 	}

// 	log.Debugf("CreateInstance, oid :%s", oid)

// 	scenario := &model.Scenario{}

// 	err := ConsumeJSON(r, scenario)

// 	if err != nil {
// 		w.WriteHeader(500)
// 		return
// 	}

// 	//etcd cid

// 	/*
// 		if ok {
// 			cid := instance.(service.Instance).ID
// 			//需要检查当前cid是否存在, 避免后台清除模块已经将container清除掉,但是session还未超时造成container无法使用
// 			_, isExits := service.CheckContainer(cid)

// 			if instance.(service.Instance).Image == scenarioConfig.Backend.Image && isExits {
// 				log.Debugf("same scenario ")
// 				ServerJSON(w, instance)
// 				return
// 			}

// 			log.Printf("Instance %s not exits ", instance.(service.Instance).ID)
// 			if isExits {
// 				service.CloseInstance(cid)
// 			}

// 			delete(session.Values, "instance")
// 			session.Save(r, w)

// 		}
// 	*/
// 	instances := []*service.Instance{}
// 	var wg sync.WaitGroup

// 	for idx := range scenario.Backend {

// 		wg.Add(1)

// 		go func() {
// 			check := service.PreCreateCheck(oid, scenario.Title, scenario.Backend[idx].Instances)

// 			if !check {
// 				wg.Done()
// 				return
// 			}

// 			instance, err := service.NewInstance(oid, scenario.Backend[idx])
// 			if err != nil {
// 				log.Error(err)
// 				http.Error(w, "Please contact service@flowq.io", 500)
// 			}
// 			log.Info(instance)
// 			instances = append(instances, instance)
// 			wg.Done()
// 		}()

// 		wg.Wait()
// 		service.UpdateInstance(oid, scenario.Title, scenario.Backend[idx].Instances)
// 	}

// 	ServerJSON(w, instances)
// }

// func ClearInstance(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)

// 	if vars["oid"] == "" || vars["cid"] == "" {
// 		http.Error(w, "Miss parameter", http.StatusBadRequest)
// 		return
// 	}

// 	status := "success"
// 	if err := service.CloseInstance(vars["oid"], vars["cid"]); err != nil {
// 		status = err.Error()
// 	}

// 	ServerJSON(w, map[string]string{"status": status, "containerID": vars["cid"]})
// }
/* ----- old func , new version use libcompose ------ */

//GetInstance get instances all/owner
func GetInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oid, ok := vars["oid"]
	if !ok || oid == "" {
		ServerJSON(w, service.GetInstance())
		return
	}

	ServerJSON(w, service.GetInstance(oid))
}

//ComposeUP create and run container from docker-compose.yml services
func ComposeUP(w http.ResponseWriter, r *http.Request) {

	log.Debugf("UP %+v", r.URL)
	vars := mux.Vars(r)
	oid, ok := vars["oid"]
	scenario := vars["sid"]

	if !ok || oid == "" {
		msg, _ := json.Marshal(map[string]string{"msg": "missing OID", "status": "400"})
		http.Error(w, string(msg), http.StatusBadRequest)
		return
	}

	log.Debugf("CreateInstance, oid :%s", oid)

	volumeReady := service.PrepareVolume(oid)

	if volumeReady == false {
		msg, _ := json.Marshal(map[string]string{"msg": "volume prepare failed", "status": "500"})
		http.Error(w, string(msg), http.StatusBadRequest)
		return
	}

	instances, err := service.ComposeUP(oid, scenario)
	if err != nil {
		msg, _ := json.Marshal(map[string]string{"msg": err.Error(), "status": "400"})
		http.Error(w, string(msg), http.StatusBadRequest)
		return
	}

	ServerJSON(w, instances)
}

//ComposeDown stop and remove instance from docker-compose.yml services
func ComposeDown(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Down %+v", r.URL)
	vars := mux.Vars(r)
	oid, ok := vars["oid"]
	scenario := vars["sid"]

	if !ok || oid == "" {
		msg, _ := json.Marshal(map[string]string{"msg": "missing OID", "status": "400"})
		http.Error(w, string(msg), http.StatusBadRequest)
		return
	}

	err := service.ComposeDown(oid, scenario)
	if err != nil {
		msg, _ := json.Marshal(map[string]string{"msg": err.Error(), "status": "400"})
		http.Error(w, string(msg), http.StatusBadRequest)
		return
	}

	ServerJSON(w, map[string]string{"msg": fmt.Sprintf(" remove %s/%s successful", oid, scenario)})
}

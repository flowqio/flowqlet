package service

import (
	"encoding/gob"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var instances = make(map[string]map[string]*Instance)

var rw = sync.Mutex{}

func ContainerExits(uid, id string) (string, bool) {
	log.Printf("%+v", instances)
	if _, ok := instances[uid]; ok {
		_, ok = instances[uid][id]
		return instances[uid][id].ScenarioName, ok
	}
	return "", false
}

func LoadInstanceFromDisk() error {
	rw.Lock()
	defer rw.Unlock()

	file, err := os.Open("instances")
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&instances)

		if err != nil {
			return err
		}
	}

	return err

}

func SaveInstanceToDisk() error {
	rw.Lock()
	defer rw.Unlock()
	file, err := os.Create("instances")
	if err == nil {
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(instances)
	}
	file.Close()
	return err
}

const instanceMaxLive = "4h"

func watch() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case _ = <-ticker.C:
			SaveInstanceToDisk()
			for _, v := range instances {

				for k1, v1 := range v {
					timeLeft := v1.ExpiresAt.Sub(time.Now())
					if timeLeft < 0 {

						//CloseInstance(k, v1.ID)
						ComposeDown(v1.OID, v1.ScenarioName)
						delete(instances[k1], k1)
						log.Infof("TimeLeft ,clear  %s instance", v1.ID)
					}
				}

			}
		}
	}
}

package service

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	InitLog()
	log.Debug("Clear Scheduler start success")
	LoadInstanceFromDisk()
	go func() {
		watch()
	}()
}

func InitLog() {
	if debug := os.Getenv("DEBUG"); debug != "" {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})
}

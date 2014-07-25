package vlog

import (
	"fmt"
	"log"
	"log/syslog"
)

var (
	Info func(string, ...interface{})
	Err  func(string, ...interface{})
)

func SetLogger(prefix string, sys bool) {
	if sys {
		logger, err := syslog.New(syslog.LOG_INFO, prefix)
		if err != nil {
			log.Fatal(err)
		}
		Info = func(f string, a ...interface{}) { logger.Info(fmt.Sprintf(f, a...)) }
		Err = func(f string, a ...interface{}) { logger.Err(fmt.Sprintf(f, a...)) }
	} else {
		log.SetPrefix(prefix)
		Info = func(f string, a ...interface{}) { log.Println(fmt.Sprintf(f, a...)) }
		Err = Info
	}
}

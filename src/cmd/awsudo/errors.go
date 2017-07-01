package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

var (
	EMissingRegion = errors.New("Missing Region")
)

func AppPanic(err error) {
	if nil != err {
		log.Warnf("Oops: %v", err)

		stackBuf := make([]byte, 2048)

		runtime.Stack(stackBuf, false)

		log.Debugf("Stack: %s", []byte(stackBuf))

		os.Exit(1)
	}
}

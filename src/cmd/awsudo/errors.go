package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	EMissingRegion = errors.New("Missing Region")
)

func AppPanic(err error) {
	if nil != err {
		log.Warnf("Oops: %v", err)

		os.Exit(1)
	}
}

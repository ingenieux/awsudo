package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func filterCurrentEnvironment() []string {
	result := make([]string, 0)

	for _, nvPair := range os.Environ() {
		nvElements := strings.SplitN(nvPair, "=", 2)

		k := nvElements[0]

		if strings.HasPrefix(k, "AWS_") {
			log.Debugf("Skipping variable '%s'", k)

			continue
		} else {
			log.Debugf("Appending variable '%s'", k)
		}

		v := ""

		if 2 == len(nvElements) {
			v = nvElements[1]
		}

		nvPairToAdd := fmt.Sprintf("%s=%s", k, v)

		result = append(result, nvPairToAdd)
	}

	return result
}

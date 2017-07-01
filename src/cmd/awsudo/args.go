package main

import (
	"bytes"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"
	"text/template"
)

var VERSION = "0.0.1-LOCAL"

const DOCOPT = `
awsudo.

Usage: awsudo [options] ROLEARN ROLESESSIONNAME [EXTERNALID]

Options:
  -h --help                 This message
  -l --logLevel=<LOGLEVEL>  Set Log Level
  -r --region=<REGION>      STS Region to Use [default: {{.Region}}]
  -s --serial=<SERIAL>      MFA Serial Number / ARN
  -t --token=<TOKENCODE>    MFA Token Code
  -v --version Version      Shows version
`

const (
	PARAM_ROLEARN         = "ROLEARN"
	PARAM_ROLESESSIONNAME = "ROLESESSIONNAME"
	PARAM_EXTERNALID      = "EXTERNALID"
	OPT_REGION            = "--region"
	OPT_LOGLEVEL          = "--logLevel"
	OPT_SERIAL            = "--serial"
	OPT_TOKEN             = "--token"
	OPT_POLICY            = "--policy"
)

func parseArguments() (map[string]interface{}, error) {
	region := "us-east-1"

	for _, envVar := range []string{"AWS_DEFAULT_REGION", "AWS_REGION"} {
		if envVarValue := os.Getenv(envVar); "" != envVarValue {
			region = envVarValue

			break
		}
	}

	templateContext := struct {
		Region string
	}{
		Region: region,
	}

	docoptTemplate := template.Must(template.New("docopt").Parse(DOCOPT))

	docoptContents := new(bytes.Buffer)

	docoptTemplate.Execute(docoptContents, templateContext)

	args, err := docopt.Parse(docoptContents.String(), nil, true, VERSION, true, false)

	if nil != err {
		log.Warnf("Oops: %v", err)
	}

	return args, err
}

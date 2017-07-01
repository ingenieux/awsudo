package main

import (
	"bytes"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"
	"text/template"
)

const DOCOPT = `awsudo.

Usage:
  awsudo [<options>] ROLEARN ROLESESSIONNAME [EXTERNALID]
  awsudo -h | --help
  awsudo -v | --version

Options:
  -h --help                 This message
  -v --version              Shows version
  -l --logLevel=<LOGLEVEL>  Set Log Level
  -r --region=<REGION>      STS Region to Use [default: {{.Region}}]
  -s --serial=<SERIAL>      MFA Serial Number / ARN
  -t --token=<TOKENCODE>    MFA Token Code
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

	args, err := docopt.Parse(docoptContents.String(), nil, true, VERSION, true, true)

	AppPanic(err)

	if logLevelToUse, ok := args[OPT_LOGLEVEL].(string); ok {
		if parsedLogLevel, err := log.ParseLevel(logLevelToUse); nil == err {
			log.SetLevel(parsedLogLevel)
		}
	}

	log.Debugf("args: %v", args)

	return args, err
}

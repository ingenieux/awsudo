package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})

	args, err := parseArguments()

	if nil != err {
		log.Fatalln(err)
	}

	sess, err := session.NewSession()

	if nil != err {
		log.Fatalln(err)
	}

	region, ok := args[OPT_REGION].(string)

	if !ok {
		log.Fatalln(EMissingRegion)
	}

	stsService := sts.New(sess, &aws.Config{
		Region: aws.String(region),
	})

	assumeRoleRequest := &sts.AssumeRoleInput{
		RoleArn:         aws.String(args[PARAM_ROLEARN].(string)),
		RoleSessionName: aws.String(args[PARAM_ROLESESSIONNAME].(string)),
	}

	if externalId, ok := args[PARAM_EXTERNALID].(string); ok {
		assumeRoleRequest.ExternalId = aws.String(externalId)
	}

	if tokenCode, ok := args[OPT_TOKEN].(string); ok {
		assumeRoleRequest.TokenCode = aws.String(tokenCode)
	}

	if serial, ok := args[OPT_SERIAL].(string); ok {
		assumeRoleRequest.SerialNumber = aws.String(serial)
	}

	if policy, ok := args[OPT_POLICY].(string); ok {
		assumeRoleRequest.Policy = aws.String(policy)
	}

	assumeRoleResults, err := stsService.AssumeRole(assumeRoleRequest)

	if nil != err {
		log.Fatalln(err)
	}

	err = executeShell(region, assumeRoleResults, args["--eval"].(bool))

	if nil != err {
		log.Fatalln(err)
	}

}

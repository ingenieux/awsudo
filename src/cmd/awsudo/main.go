package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})

	args, err := parseArguments()

	AppPanic(err)

	sess, err := session.NewSession()

	AppPanic(err)

	region, ok := args[OPT_REGION].(string)

	if !ok {
		AppPanic(EMissingRegion)
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

	AppPanic(err)

	log.Debugf("AssumeRole: %+v", *assumeRoleResults)

	environ := filterCurrentEnvironment()

	unixExpiration := assumeRoleResults.Credentials.Expiration.Unix()

	for k, v := range map[string]string{
		"AWS_ACCESS_KEY_ID":             *assumeRoleResults.Credentials.AccessKeyId,
		"AWS_SECRET_ACCESS_KEY":         *assumeRoleResults.Credentials.SecretAccessKey,
		"AWS_DEFAULT_REGION":            region,
		"AWS_SESSION_TOKEN":             *assumeRoleResults.Credentials.SessionToken,
		"AWS_SESSION_ROLE_USER_ARN":     *assumeRoleResults.AssumedRoleUser.Arn,
		"AWS_SESSION_EXPIRATION":        fmt.Sprintf("%d", unixExpiration),
		"AWS_SESSION_ROLE_USER_ROLE_ID": *assumeRoleResults.AssumedRoleUser.AssumedRoleId,
	} {
		if v != "" {
			log.Debugf("Appending variable '%s'", k)

			environ = append(environ, k+"="+v)
		}
	}

	environ = append(environ)

	shell := os.Getenv("SHELL")

	shellCmd := exec.Command(shell)

	shellCmd.Stderr = os.Stderr
	shellCmd.Stdout = os.Stdout
	shellCmd.Stdin = os.Stdin

	shellCmd.Env = environ

	shellCmd.Run()
}

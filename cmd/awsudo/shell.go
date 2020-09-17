package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sts"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func executeShell(region string, stsRole *sts.AssumeRoleOutput, evalMode bool) error {
	log.Debugf("AssumeRole: %+v", *stsRole)

	unixExpiration := stsRole.Credentials.Expiration.Unix()

	envVars := map[string]string{
		"AWS_ACCESS_KEY_ID":             *stsRole.Credentials.AccessKeyId,
		"AWS_SECRET_ACCESS_KEY":         *stsRole.Credentials.SecretAccessKey,
		"AWS_DEFAULT_REGION":            region,
		"AWS_SESSION_TOKEN":             *stsRole.Credentials.SessionToken,
		"AWS_SESSION_ROLE_USER_ARN":     *stsRole.AssumedRoleUser.Arn,
		"AWS_SESSION_EXPIRATION":        fmt.Sprintf("%d", unixExpiration),
		"AWS_SESSION_ROLE_USER_ROLE_ID": *stsRole.AssumedRoleUser.AssumedRoleId,
	}

	if evalMode {
		var statements []string

		for k, v := range envVars {
			statements = append(statements, fmt.Sprintf("%s=%s", k, v))
		}

		fmt.Println("export " + strings.Join(statements, " "))

		return nil
	}

	environ := filterCurrentEnvironment()

	for k, v := range envVars {
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

	return shellCmd.Run()
}

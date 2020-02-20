package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)
const (
	DEFAULTREGION = "us-east-1"
)
type AwsDefaults struct {
	region string
}

func (a AwsDefaults) Session() *session.Session {

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(DEFAULTREGION),
		LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody)},
	)
	return sess
}

func Session() *session.Session {
	return AwsDefaults{region:DEFAULTREGION}.Session()
}
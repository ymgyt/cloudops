package backends

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/ymgyt/cloudops/core"
)

// NewAWSSession -
func NewAWSSession(region, id, secret, token string) (*session.Session, error) {
	crd := credentials.NewStaticCredentials(id, secret, token)
	cfg := aws.NewConfig().WithCredentials(crd).WithRegion(region)
	return session.NewSession(cfg)
}

// WrapAWSError -
func WrapAWSError(err awserr.Error) error {
	var w error
	switch err.Code() {
	case "EmptyStaticCreds":
		w = core.NewError(core.Unauthenticated, "empty aws credentials check AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY env")
	default:
		w = core.WrapError(core.Internal, "", err)
	}
	return w
}

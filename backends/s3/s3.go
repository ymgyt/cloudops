package s3

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"

	"github.com/ymgyt/cloudops/backends"
	"github.com/ymgyt/cloudops/core"
)

// NewS3
func New(ctx *core.Context, region, accessKeyID, secretAccessKey, token string) (core.Backend, error) {
	return newS3(ctx, region, accessKeyID, secretAccessKey, token)
}

func newS3(ctx *core.Context, region, accessKeyID, secretAccessKey, token string) (*s3Client, error) {
	sess, err := backends.NewAWSSession(region, accessKeyID, secretAccessKey, token)
	if err != nil {
		// FIXME mote appropriate code
		return nil, core.WrapError(core.Internal, "backend.NewAWSSession", err)
	}
	return &s3Client{ctx: ctx, client: s3.New(sess)}, nil
}

type s3Client struct {
	ctx    *core.Context
	client *s3.S3
}

// Put -
func (c *s3Client) Put(in *core.PutInput) (*core.PutOutput, error) {
	var n int
	for _, r := range in.Resources {
		err := c.put(r, in)
		if err != nil {
			return nil, c.wrapError(err)
		}
		n++
	}
	return &core.PutOutput{PutNum: n}, nil
}

// Fetch -
func (c *s3Client) Fetch(in *core.FetchInput) (*core.FetchOutput, error) {
	return nil, core.NotImplementedError("s3Client.Fetch()")
}

// Remove -
func (c *s3Client) Remove(in *core.RemoveInput) (*core.RemoveOutput, error) {
	return nil, core.NotImplementedError("s3Client.Remove()")
}

func (c *s3Client) put(r core.Resource, in *core.PutInput) error {
	input, err := c.convertToPutInput(r, in)
	if err != nil {
		return err
	}
	defer func() {
		if input.Body == nil {
			return
		}
		if clsErr := input.Body.(io.Closer).Close(); err == nil {
			err = clsErr
		}
	}()
	log := c.ctx.Log.With(zap.String("src", r.URI()), zap.String("dest", fmt.Sprintf("s3://%s/%s", *input.Bucket, *input.Key)))
	if in.Dryrun {
		log.Info("copy", zap.Bool("dryrun", in.Dryrun))
		return nil
	}
	log.Info("copy")
	_, err = c.client.PutObject(input)
	return err
}

func (c *s3Client) convertToPutInput(r core.Resource, in *core.PutInput) (*s3.PutObjectInput, error) {
	bucket, key, err := c.split(in.Dest)
	if err != nil {
		return nil, err
	}
	if in.Recursive {
		key = path.Join(key, path.Base(r.URI()))
	}
	body, err := r.Open()
	if err != nil {
		return nil, err
	}
	// FIXME: um... should core.Resource interface change ?
	// reading all to buffer, then converting to io.Seeker is ng.
	seeker, ok := body.(io.ReadSeeker)
	if !ok {
		return nil, core.NewError(core.Internal, ("currently require resource to implement io.Seeker"))
	}
	return &s3.PutObjectInput{
		Body:   seeker,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, nil
}

func (c *s3Client) split(path string) (bucket, key string, err error) {
	org := path
	fail := func() (string, string, error) {
		return "", "", core.NewError(core.InvalidParam, fmt.Sprintf("invalid s3 path %s, should be s3://<bucket>/<object> format", org))
	}
	if !strings.HasPrefix(path, "s3://") {
		return fail()
	}
	path = path[5:]
	idx := strings.Index(path, "/")
	if idx == -1 {
		return fail()
	}

	bucket = path[:idx]
	if idx == len(path) {
		return fail()
	}
	key = path[idx+1:]
	if bucket == "" || key == "" {
		return fail()
	}

	return bucket, key, nil
}

func (c *s3Client) wrapError(err error) error {
	if err == nil {
		return nil
	}
	if coreErr, ok := err.(core.Error); ok {
		return coreErr
	}
	if awsErr, ok := err.(awserr.Error); ok {
		return backends.WrapAWSError(awsErr)
	}
	return core.WrapError(core.Internal, "s3Client.Put", err)
}

type s3Resource struct {
	path string
}

// Type -
func (r *s3Resource) Type() core.ResourceType {
	return core.S3Resource
}

// URI -
func (r *s3Resource) URI() string {
	return r.path
}

// Open -
func (r *s3Resource) Open() (io.ReadCloser, error) {
	return nil, core.NewError(core.NotImplementedYet, "s3Resource.Open()")
}

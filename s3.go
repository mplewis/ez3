package ez3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// S3Client is the interface for an AWS S3-compatible client.
type S3Client interface {
	PutObject(context.Context, *awsS3.PutObjectInput, ...func(*awsS3.Options)) (*awsS3.PutObjectOutput, error)
	GetObject(context.Context, *awsS3.GetObjectInput, ...func(*awsS3.Options)) (*awsS3.GetObjectOutput, error)
	DeleteObject(context.Context, *awsS3.DeleteObjectInput, ...func(*awsS3.Options)) (*awsS3.DeleteObjectOutput, error)
	ListObjectsV2(context.Context, *awsS3.ListObjectsV2Input, ...func(*awsS3.Options)) (*awsS3.ListObjectsV2Output, error)
}

// S3EZ3 is an implementation of EZ3 backed by an S3-compatible file store.
type S3EZ3 struct {
	bucket    string
	namespace string
	client    S3Client
}

// S3Args is the set of arguments for creating a new S3-backed EZ3.
type S3Args struct {
	Bucket    string   // Required. The bucket that holds stored data.
	Namespace string   // Required. The namespace for this instance's keys.
	Client    S3Client // Optional. If not provided, autoconfigures an AWS S3 client from your environment.
}

// S3ClientArgs is the set of arguments for creating a new S3 client for use with S3EZ3.
type S3ClientArgs struct {
	Endpoint     string // Required. The web endpoint of the S3 service. Usually starts with https://
	Region       string // Required. The region of the S3 service.
	UsePathStyle bool   // Optional. If true, the S3 client will use path-style addressing.
}

// notFoundErr generates a custom error for a missing key.
func notFoundErr(key string) error {
	return fmt.Errorf("key not found: %s", key)
}

// wasNotFound returns true if the error represents an S3 "not found" error.
func wasNotFound(err error) bool {
	if err == nil {
		return false
	}
	var ae smithy.APIError
	if !errors.As(err, &ae) {
		return false
	}
	return ae.ErrorCode() == "NoSuchKey"
}

// ns adds the namespace to a given key.
func (s *S3EZ3) ns(key string) string {
	return s.namespace + "/" + key
}

// Get retrieves a value from S3.
func (s S3EZ3) Get(key string, dst Serializable) error {
	output, err := s.client.GetObject(context.TODO(), &awsS3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.ns(key)),
	})

	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) && ae.ErrorCode() == "NoSuchKey" {
			return notFoundErr(key)
		}
		return err
	}

	data, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return err
	}
	return dst.Deserialize(data)
}

// Set stores a value in S3.
func (s S3EZ3) Set(key string, val Serializable) error {
	data, err := val.Serialize()
	if err != nil {
		return err
	}
	_, err = s.client.PutObject(context.TODO(), &awsS3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.ns(key)),
		Body:   bytes.NewReader(data),
	})
	return err
}

// Del removes a value from S3.
func (s S3EZ3) Del(key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &awsS3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.ns(key)),
	})
	if wasNotFound(err) {
		return nil
	}
	return err
}

// List lists all keys in the namespace with the given prefix.
func (s S3EZ3) List(prefix string) ([]string, error) {
	// TODO: Paginate
	output, err := s.client.ListObjectsV2(context.TODO(), &awsS3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(s.ns(prefix)),
	})
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, object := range output.Contents {
		fullKey := aws.ToString(object.Key)
		key := strings.TrimPrefix(fullKey, s.namespace+"/")
		keys = append(keys, key)
	}
	return keys, nil
}

// NewS3 creates a new S3-based EZ3 client.
func NewS3(args S3Args) (EZ3, error) {
	if args.Bucket == "" {
		return nil, errors.New("bucket not specified")
	}
	if args.Namespace == "" {
		return nil, errors.New("namespace not specified")
	}
	if args.Client == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		args.Client = awsS3.NewFromConfig(cfg)
	}
	return S3EZ3{bucket: args.Bucket, namespace: args.Namespace, client: args.Client}, nil
}

// NewS3Client creates a new S3 client for use with S3EZ3.
// This is a helper function for users who want to connect to an S3-compatible cloud that is not AWS.
func NewS3Client(args S3ClientArgs) (*awsS3.Client, error) {
	if args.Endpoint == "" {
		return nil, errors.New("endpoint not specified")
	}
	if args.Region == "" {
		return nil, errors.New("region not specified")
	}
	if !strings.HasPrefix(args.Endpoint, "https://") {
		fmt.Printf("WARNING: S3 endpoint %s does not start with https://. "+
			"This is usually a mistake.\n", args.Endpoint)
	}
	resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service != awsS3.ServiceID {
				return aws.Endpoint{}, &aws.EndpointNotFoundError{} // fallback
			}
			return aws.Endpoint{URL: args.Endpoint, SigningRegion: args.Region}, nil
		})
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(resolver))
	if err != nil {
		return nil, err
	}
	return awsS3.NewFromConfig(cfg, func(o *awsS3.Options) {
		o.UsePathStyle = args.UsePathStyle
	}), nil
}

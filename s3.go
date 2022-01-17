package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type s3ez3 struct {
	bucket    string
	namespace string
	client    *s3.Client
}

func (s *s3ez3) ns(key string) string {
	return s.namespace + "/" + key
}

func (s s3ez3) Get(key string, dst Serdeable) error {
	output, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.ns(key)),
	})

	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) && ae.ErrorCode() == "NoSuchKey" {
			return fmt.Errorf("key not found: %s", key)
		}
		return err
	}

	data, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return err
	}
	return dst.Deserialize(data)
}

func (s s3ez3) Set(key string, val Serdeable) error {
	data, err := val.Serialize()
	if err != nil {
		return err
	}
	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.ns(key)),
		Body:   bytes.NewReader(data),
	})
	return err
}

func (s s3ez3) Del(key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.ns(key)),
	})
	return err
}

func (s s3ez3) List(prefix string) ([]string, error) {
	output, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(s.ns(prefix)),
	})
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, object := range output.Contents {
		// TODO: Strip prefix
		keys = append(keys, aws.ToString(object.Key))
	}
	return keys, nil
}

// NewS3 creates a new S3-based EZ3 client.
func NewS3() (EZ3, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := s3.NewFromConfig(cfg)

	s := s3ez3{
		client:    client,
		bucket:    "mplewis-s3kv-test",
		namespace: "s3ez3",
	}
	return s, nil
}

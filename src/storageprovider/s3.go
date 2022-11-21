package storageprovider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	appconfig "github.com/sid-sun/ioctl-api/config"
)

type S3Provider struct {
	bucketname string
	client     *s3.Client
	acl        types.ObjectCannedACL
}

func InitS3StorageProvider() *S3Provider {
	var cfg aws.Config
	var acl types.ObjectCannedACL

	// Dynamically Load config for S3 Provider
	switch appconfig.Cfg.S3.Provider {
	case "R2":
		cfg = R2()
	case "S3":
		cfg = S3()
		acl = types.ObjectCannedACLPublicRead
	}

	client := s3.NewFromConfig(cfg)

	return &S3Provider{
		client:     client,
		bucketname: appconfig.Cfg.S3.Bucket,
		acl:        acl,
	}
}

func (sp *S3Provider) UploadSnippet(data io.Reader, id string) error {
	_, err := sp.GetObjectInfo(id)
	if err != nil && err != ErrNotFound {
		return err
	}

	if err != ErrNotFound {
		return ErrAlreadyExists
	}

	t := time.Now().Add(time.Hour)
	_, err = sp.client.PutObject(context.Background(), &s3.PutObjectInput{
		Key:     &id,
		Body:    data,
		Bucket:  &sp.bucketname,
		ACL:     sp.acl,
		Expires: &t,
	})
	if err != nil {
		return err
	}

	return nil
}

func (sp *S3Provider) GetObjectInfo(id string) (*s3.HeadObjectOutput, error) {
	objhead, err := sp.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: &sp.bucketname,
		Key:    &id,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			var nf types.NotFound
			switch apiErr.ErrorCode() {
			case nf.ErrorCode():
				return nil, ErrNotFound
			}
		}
		fmt.Println("'lost it '")
		return nil, err
	}
	return objhead, nil
}

func (sp *S3Provider) DownloadSnippet(id string) (data []byte, err error) {
	obj, err := sp.client.GetObject(context.Background(), &s3.GetObjectInput{
		Key:    &id,
		Bucket: &sp.bucketname,
	})
	// TODO: Handle not found error
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	objData, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}

	return objData, nil
}

package storageprovider

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/fitant/xbin-api/config"
)

type S3Provider struct {
	service    *s3.S3
	bucketname string
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	acl        string
}

func InitS3StorageProvider() *S3Provider {
	sess := session.Must(session.NewSession())

	service := s3.New(sess, aws.NewConfig())
	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)

	return &S3Provider{
		service:    service,
		uploader:   uploader,
		downloader: downloader,
		bucketname: config.Cfg.S3.Bucket,
		acl:        "public-read",
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

	_, err = sp.uploader.Upload(&s3manager.UploadInput{
		Key:    &id,
		Body:   data,
		Bucket: &sp.bucketname,
		ACL:    &sp.acl,
	})
	if err != nil {
		return err
	}
	return nil
}

func (sp *S3Provider) GetObjectInfo(id string) (*s3.HeadObjectOutput, error) {
	objhead, err := sp.service.HeadObject(&s3.HeadObjectInput{
		Bucket: &sp.bucketname,
		Key:    &id,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				return nil, ErrNotFound
			default:
				return nil, err
			}
		}
		return nil, err
	}
	return objhead, nil
}

func (sp *S3Provider) DownloadSnippet(id string) (data []byte, err error) {
	objhead, err := sp.GetObjectInfo(id)
	if err != nil {
		return nil, err
	}

	data = make([]byte, *objhead.ContentLength)

	_, err = sp.downloader.Download(aws.NewWriteAtBuffer(data), &s3.GetObjectInput{
		Key:    &id,
		Bucket: &sp.bucketname,
	})
	if err != nil {
		return nil, err
	}

	return
}

package cloud

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"os"
)

// S3StoreConfig is the S3 store configuration.
type S3StoreConfig struct {
	// AccessKey is the access key.
	AccessKey string
	// SecretKey is the secret key.
	SecretKey string
	// Endpoint is the endpoint.
	Endpoint string
	// Bucket is the bucket.
	Bucket string
	// Path is the path.
	Path string
}

// S3Store is the interface of S3 data store.
// It is used to store real-time batteries data in S3.
// The real-time batteries data is used to do e.
type S3Store interface {
	// Init initializes the S3 store.
	Init() error

	// Upload uploads the file to S3.
	Upload(fileName string) error

	// Download downloads the file from S3.
	Download(fileName string) error
}

type s3StoreImpl struct {
	cfg *S3StoreConfig
	svc *s3.S3
}

// NewS3Store creates a new S3 store.
func NewS3Store(cfg *S3StoreConfig) S3Store {
	return &s3StoreImpl{cfg: cfg}
}

// Init initializes the S3 store.
func (s *s3StoreImpl) Init() error {
	// create S3 client
	creds := credentials.NewStaticCredentials(s.cfg.AccessKey, s.cfg.SecretKey, "")
	_, err := creds.Get()
	if err != nil {
		return fmt.Errorf("failed to get credentials: %w", err)
	}

	awsConfig := &aws.Config{
		Credentials: creds,
		Endpoint:    aws.String(s.cfg.Endpoint),
		Region:      aws.String(s.cfg.Region),
	}

	request.WithRetryer(awsConfig, aws.NewDefaultRetryer())

	svc := s3.New(session.New(), awsConfig)

	// create bucket if not exists
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(s.cfg.Bucket),
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	s.svc = svc

	return nil
}

// Upload uploads the file to S3.
func (s *s3StoreImpl) Upload(fileName string) error {
	if s.svc == nil {
		return errors.New("S3 store is not initialized")
	}

	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = s.svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(s.cfg.Path + "/" + fileName),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// Download downloads the file from S3.
func (s *s3StoreImpl) Download(fileName string) error {
	if s.svc == nil {
		return errors.New("S3 store is not initialized")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = s.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(s.cfg.Path + "/" + fileName),
	})
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	return nil
}

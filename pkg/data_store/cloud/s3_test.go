package cloud

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/s3/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockS3Client is a mock implementation of the S3 client.
type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) CreateBucket(params *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	args := m.Called(params)
	return args.Get(0).(*s3.CreateBucketOutput), args.Error(1)
}

func (m *MockS3Client) PutObject(params *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := m.Called(params)
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

func TestS3StoreInit(t *testing.T) {
	// Create a mock S3 client
	mockS3Client := &MockS3Client{}

	// Create a test S3StoreImpl with the mock client
	s3Store := &s3StoreImpl{
		cfg: &S3StoreConfig{
			AccessKey: "test-access-key",
			SecretKey: "test-secret-key",
			Endpoint:  "http://localhost:9000", // Replace with your S3 endpoint
			Bucket:    "test-bucket",
			Path:      "test-path",
		},
		svc: mockS3Client,
	}

	// Configure expectations for CreateBucket method
	mockS3Client.On("CreateBucket", mock.Anything).Return(&s3.CreateBucketOutput{}, nil)

	// Call the Init method
	err := s3Store.Init()

	// Assert that the expectations were met
	mockS3Client.AssertExpectations(t)

	// Check for errors
	assert.Nil(t, err, "Init should not return an error")
}

func TestS3StoreUpload(t *testing.T) {
	// Create a mock S3 client
	mockS3Client := &MockS3Client{}

	// Create a test S3StoreImpl with the mock client
	s3Store := &s3StoreImpl{
		cfg: &S3StoreConfig{
			AccessKey: "test-access-key",
			SecretKey: "test-secret-key",
			Endpoint:  "http://localhost:9000", // Replace with your S3 endpoint
			Bucket:    "test-bucket",
			Path:      "test-path",
		},
		svc: mockS3Client,
	}

	// Configure expectations for PutObject method
	mockS3Client.On("PutObject", mock.Anything).Return(&s3.PutObjectOutput{}, nil)

	// Call the Upload method
	err := s3Store.Upload("test-file.txt")

	// Assert that the expectations were met
	mockS3Client.AssertExpectations(t)

	// Check for errors
	assert.Nil(t, err, "Upload should not return an error")
}

func TestS3StoreDownload(t *testing.T) {
	// Create a mock S3 client
	mockS3Client := &MockS3Client{}

	// Create a test S3StoreImpl with the mock client
	s3Store := &s3StoreImpl{
		cfg: &S3StoreConfig{
			AccessKey: "test-access-key",
			SecretKey: "test-secret-key",
			Endpoint:  "http://localhost:9000", // Replace with your S3 endpoint
			Bucket:    "test-bucket",
			Path:      "test-path",
		},
		svc: mockS3Client,
	}

	// Configure expectations for GetObject method
	mockS3Client.On("GetObject", mock.Anything).Return(&s3.GetObjectOutput{}, nil)

	// Call the Download method
	err := s3Store.Download("test-file.txt")

	// Assert that the expectations were met
	mockS3Client.AssertExpectations(t)

	// Check for errors
	assert.Nil(t, err, "Download should not return an error")
}

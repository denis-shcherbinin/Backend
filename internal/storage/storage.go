package storage

import (
	"errors"
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"math/rand"
	"strings"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fileNameLength = 16
)

type Storage struct {
	client   *s3.S3
	bucket   string
	endpoint string
}

type UploadInput struct {
	Body        string
	ContentType string
}

func NewStorage(cfg *config.StorageConfig) (*Storage, error) {
	key := cfg.AccessKey
	secret := cfg.SecretKey

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(cfg.Endpoint),
		Region:      aws.String(cfg.Region),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, err
	}

	return &Storage{
		client:   s3.New(newSession),
		bucket:   cfg.Name,
		endpoint: cfg.Endpoint,
	}, nil
}

func (s *Storage) Upload(input UploadInput) (string, error) {
	fileName := s.generateFileName()

	object := s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(fileName),
		Body:        strings.NewReader(input.Body),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(input.ContentType),
	}

	_, err := s.client.PutObject(&object)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Storage: object: %v upload error", object))
	}

	return s.generateFileURL(fileName), nil
}

func (s *Storage) Delete(imageURL string) error {
	fileName := imageURL[len(imageURL)-fileNameLength:]
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
	}

	_, err := s.client.DeleteObject(input)

	return err
}

func (s *Storage) generateFileName() string {
	b := make([]byte, fileNameLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func (s *Storage) generateFileURL(fileName string) string {
	return fmt.Sprintf("https://%s/%s/%s", s.endpoint, s.bucket, fileName)
}

package backend

import (
	"appstore/constants"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"mime/multipart"
)

var (
	S3Backend *S3StorageBackend
)
 
 
type S3StorageBackend struct {
	s3Session *s3.S3
	uploader *s3manager.Uploader
}

func InitS3Backend() {
	
	s := s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(constants.S3_REGION_NAME),
		Credentials: credentials.NewStaticCredentials(constants.S3_BUCKET_KEY_ID, constants.S3_BUCKET_KEY, ""),
    })))
	
	_, err := s.Config.Credentials.Get();

	if err!= nil {
        panic(err)
    }

	if s == nil {
		fmt.Println("Failed to create S3 session")
        return
    }

	uploader := s3manager.NewUploaderWithClient(s)
	S3Backend = &S3StorageBackend{
        s3Session: s,
		uploader: uploader,
    }
}

func (s *S3StorageBackend) SaveToS3(file multipart.File, id string) (string, error) {
	
	filename := (uuid.New()).String()
	params := &s3.PutObjectInput{
		Bucket: aws.String(constants.S3_BUCKET_NAME),
		Key:    aws.String(filename),
		ACL:    aws.String("public-read"),
		Body:   file,
		ContentType: aws.String("image/jpeg"),
	   }
	response, err := s.s3Session.PutObject(params)
	if err!= nil {
		fmt.Println("Failed to upload image with S3 session")
		return "", err
	}	
	fmt.Println("Successfully uploaded image with S3 session")
	fmt.Println(response)
	fmt.Println("https://" + constants.S3_BUCKET_NAME + ".s3.amazonaws.com/" + filename)
	return "https://" + constants.S3_BUCKET_NAME + ".s3.amazonaws.com/" + filename, nil
}


 
 
 
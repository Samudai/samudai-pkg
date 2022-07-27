package fileupload

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	uuid "github.com/satori/go.uuid"
)

var (
	s3Client          *s3.S3
	endpoint          = os.Getenv("ENDPOINT")
	bucketName        = os.Getenv("BUCKET_NAME")
	allowedExtensions []string
	allowedSize 	 int64
)

func Init(allowedExt []string, allowSize int64) {
	allowedExtensions = allowedExt
	allowedSize = allowSize
	key := os.Getenv("SPACES_KEY")
	secret := os.Getenv("SPACES_SECRET")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String("https://" + endpoint),
		Region:      aws.String("us-east-1"),
	}

	newSession := session.New(s3Config)
	s3Client = s3.New(newSession)
}

func UploadFile(file *multipart.FileHeader, filename string) (string, error) {
	var url string
	buffer := make([]byte, file.Size)
	f, err := file.Open()
	if err != nil {
		return url, err
	}
	defer f.Close()
	f.Read(buffer)

	ext := filepath.Ext(file.Filename)
	if !isAllowedExtension(ext) {
		return url, fmt.Errorf("%s is not an allowed file extension", ext)
	}
	if file.Size > allowedSize {
		return url, fmt.Errorf("%s is too large", file.Filename)
	}
	fileName := uuid.NewV4().String() + ext

	object := s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(fileName),
		Body:          bytes.NewReader(buffer),
		ACL:           aws.String("public-read"),
		ContentLength: aws.Int64(file.Size),
		ContentType:   aws.String(http.DetectContentType(buffer)),
	}

	_, err = s3Client.PutObject(&object)
	if err != nil {
		return url, err
	}
	url = "https://" + bucketName + "." + endpoint + "/" + fileName

	return url, nil
}

func isAllowedExtension(ext string) bool {
	for _, allowed := range allowedExtensions {
		if strings.ToLower(allowed) == strings.ToLower(ext) {
			return true
		}
	}
	return false
}

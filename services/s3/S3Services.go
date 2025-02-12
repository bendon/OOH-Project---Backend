package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3ClientCon *s3.Client
)

func InitializeS3Client() {
	// ctx := context.TODO()

	// Load environment variables
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	region := os.Getenv("S3_REGION")
	// Scaleway region

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v", err)
	}

	// Example: List Objects in Bucket
	// listObjects(ctx, s3Client, bucket)
	s3ClientCon = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = region
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(endpoint)
		o.EndpointResolverV2 = s3.NewDefaultEndpointResolverV2()
		o.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			}, nil
		})
	})
}

func ListObjects(ctx context.Context, client *s3.Client, bucket string) {
	resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		fmt.Printf("Failed to list objects in bucket %s: %v", bucket, err)
	}

	for _, obj := range resp.Contents {
		fmt.Println(*obj.Key)
	}
}

func UploadFileToBucket(key string, file multipart.File, fileHeader *multipart.FileHeader) error {
	if s3ClientCon == nil {
		// Initialize S3 client
		InitializeS3Client()
	}

	ctx := context.TODO()
	bucket := os.Getenv("S3_BUCKET")

	// Read file into memory buffer
	fileBuffer := new(bytes.Buffer)
	_, err := io.Copy(fileBuffer, file)
	if err != nil {
		fmt.Println("Failed to read file into memory buffer:", err)
		return err
	}

	// Convert buffer to io.ReadSeeker (which supports rewinding)
	fileReader := bytes.NewReader(fileBuffer.Bytes())

	// Get content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream" // Default fallback
	}
	//log the url endpoint

	// Upload file to S3
	_, err = s3ClientCon.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        fileReader,
		ContentType: &contentType,
	})
	if err != nil {
		fmt.Println("Failed to upload file to S3:", err)
		return err
	}
	return nil

}

func DownloadFileFromBucket(key string) (io.ReadCloser, error) {
	if s3ClientCon == nil {
		// Initialize S3 client
		InitializeS3Client()
	}
	ctx := context.TODO()
	bucket := os.Getenv("S3_BUCKET")

	// Download file from S3
	resp, err := s3ClientCon.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		fmt.Println("Failed to download file:", err)
	}
	return resp.Body, nil

}

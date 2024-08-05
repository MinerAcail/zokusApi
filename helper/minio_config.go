package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	// "github.com/minio/minio-go/v7/pkg/cors"
)

func InitMinIO(endpoint, accessKey, secretKey string) (*minio.Client, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // Set to true if using HTTPS
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %v", err)
	}
	
	return minioClient, nil
}

func CreateBucketInMinio(ctx context.Context, minioClient *minio.Client, id string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	
	// Generate a unique bucket name
	bucketName := fmt.Sprintf("vendor-%s-bucket", id)

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %v", err)
	}
	if !exists {
		if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
	}

	// Set the bucket policy to public
	policyJSON := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucketName)
	err = minioClient.SetBucketPolicy(ctx, bucketName, policyJSON)
	if err != nil {
		return fmt.Errorf("could not set bucket policy: %v", err)
	}

	

	return nil
}

func GenerateFileURL(bucketName, fileName string, minioClient *minio.Client) (string, error) {
	expiration := time.Hour * 24 // URL expiration time (24 hours)

	// Generate the presigned URL with a duration for expiration
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, fileName, expiration, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}
	return presignedURL.String(), nil
}

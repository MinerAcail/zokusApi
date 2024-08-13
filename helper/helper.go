package helper

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/mineracail/zokus/graph/model"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

var IsTrue bool = true
var IsFalse bool = false

func ReadFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read the sql file %s %w", filePath, err)
	}
	return content, nil
}

func DeleteByID(db *gorm.DB, model interface{}, id string) error {
	// Find the model instance by ID
	if err := db.Where("id = ?", id).First(model).Error; err != nil {
		return fmt.Errorf("%s not found: %v", reflect.TypeOf(model), err)
	}

	// Delete the model instance from the database
	if err := db.Delete(model).Error; err != nil {
		return fmt.Errorf("failed to delete %s: %v", reflect.TypeOf(model), err)
	}

	return nil
}

func GetAllByModel(db *gorm.DB, model interface{}, preloads ...string) error {
	query := db.Order("created_at DESC")

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Find(model).Error; err != nil {
		return fmt.Errorf("failed to get response for %s: %v", reflect.TypeOf(model), err)
	}

	return nil
}

func CheckDuplicateRecords(db *gorm.DB, model interface{}, inputName, inputPhoneNumber string) error {
	// Check if a user with the same name and number already exists
	err := db.Where("name = ? AND number = ?", inputName, inputPhoneNumber).First(model).Error
	if err == nil {
		return fmt.Errorf(" A member with the number already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to query for duplicates: %w", err)
	}

	// Check if a user with the same number already exists
	err = db.Where("number = ?", inputPhoneNumber).First(model).Error
	if err == nil {
		return fmt.Errorf(" A member with the same number already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to query for duplicates: %w", err)
	}

	return nil
}

func CleanPhoneNumber(phoneNumber string) string {
	cleanedPhoneNumber := strings.ReplaceAll(phoneNumber, " ", "")
	cleanedPhoneNumber = strings.TrimPrefix(cleanedPhoneNumber, "+233")
	cleanedPhoneNumber = strings.TrimPrefix(cleanedPhoneNumber, "233")

	// Check if the cleaned phone number has fewer than 10 digits
	if len(cleanedPhoneNumber) < 10 {
		// Add a leading "0" to make it 10 digits
		cleanedPhoneNumber = "0" + cleanedPhoneNumber
	}

	return cleanedPhoneNumber
}

func FetchUser(db *gorm.DB, userID uuid.UUID) (*model.User, error) {
	var user model.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return &user, nil
}

func FetchProducts(db *gorm.DB, productIDs []uuid.UUID) ([]*model.Product, error) {
	var products []*model.Product
	if err := db.Where("id IN (?)", productIDs).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("some products not found: %v", err)
	}
	return products, nil
}

func UploadFileToMinio(ctx context.Context, minioClient *minio.Client, bucketName, fileType string, file *graphql.Upload) (string, error) {
	uniqueID := uuid.New().String()
	fileName := fmt.Sprintf("%s-%s-%s", bucketName, fileType, uniqueID)

	// Ensure bucket exists before uploading
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to check if bucket exists: %v", err)
	}
	if !exists {
		return "", fmt.Errorf("bucket %s does not exist", bucketName)
	}

	// Upload the file to MinIO
	_, err = minioClient.PutObject(ctx, bucketName, fileName, file.File, file.Size, minio.PutObjectOptions{ContentType: file.ContentType})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %v", err)
	}

	// Generate and return the file URL
	fileURL := fmt.Sprintf("/minio/%s/%s", bucketName, fileName)

	return fileURL, nil
}

// Helper function to determine user type
func DetermineUserType(db *gorm.DB, userUUID uuid.UUID) model.UserType {
	var user model.User
	if err := db.Where("id = ?", userUUID).First(&user).Error; err == nil {
		return model.UserTypeUser
	}

	var vendor model.Vendor
	if err := db.Where("id = ?", userUUID).First(&vendor).Error; err == nil {
		return model.UserTypeVendor
	}

	var sideSwapper model.SideSwapper
	if err := db.Where("id = ?", userUUID).First(&sideSwapper).Error; err == nil {
		return model.UserTypeSideswapper
	}

	return ""
}

// fileURL := fmt.Sprintf("/minio/%s/%s", bucketName, fileName)
// bucketName := fmt.Sprintf("vendor-%s-bucket", vendorIDStr)

func ImageFromMinIOURL(ctx context.Context, minioClient *minio.Client, fileURL string) ([]byte, error) {
	// Extract bucketName and fileName from the URL
	parts := strings.SplitN(fileURL, "/", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid file URL format: %s", fileURL)
	}
	bucketName := parts[2]
	fileName := parts[3]

	object, err := minioClient.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("error fetching file from MinIO: %w", err)
	}
	defer object.Close()

	imageData, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("error reading file from MinIO: %w", err)
	}
	return imageData, nil
}

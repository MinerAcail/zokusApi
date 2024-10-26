package helper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"

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

// deleteNestedComments deletes all replies for the given comments recursively
func DeleteNestedComments(db *gorm.DB, comments []*model.Comment) error {
	for _, reply := range comments {
		// Fetch the reply to get its nested replies
		var nestedReplies []*model.Comment
		if err := db.Preload("Replies").Find(&nestedReplies, "parent_comment_id = ?", reply.ID).Error; err != nil {
			return fmt.Errorf("failed to fetch nested replies: %w", err)
		}

		// Recursive call to delete nested replies
		if err := DeleteNestedComments(db, nestedReplies); err != nil {
			return err
		}

		// Delete the reply itself
		if err := db.Delete(reply).Error; err != nil {
			return fmt.Errorf("failed to delete reply: %w", err)
		}
	}
	return nil
}

func GenerateOrderCartID(cartID string) (string, error) {
	// Create a new random generator seeded with the current time
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate three random lowercase letters
	randomLetters := fmt.Sprintf("%c%c%c", 'a'+r.Intn(26), 'a'+r.Intn(26), 'a'+r.Intn(26))
	prefix := cartID[:3]
	if len(cartID) < 3 {
		return "", errors.New("cartID must be at least 3 characters long")
	}

	// Return the formatted string
	return fmt.Sprintf("%s%s", prefix, randomLetters), nil
}

func CreateMetaData(db *gorm.DB, cartItemID uuid.UUID, cartID uuid.UUID) error {
	// Check if the purchase meta data already exists
	var existingPurchaseMetaData model.PurchaseMetaData
	if err := db.Where("cart_item_id = ? AND cart_id = ?", cartItemID, cartID).First(&existingPurchaseMetaData).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			// Record doesn't exist, create a new one
			purchaseMetaData := &model.PurchaseMetaData{
				CartItemID: cartItemID,
				Status:     string(model.StatusInCart),
				CartID:     cartID,
			}
			if err := db.Create(&purchaseMetaData).Error; err != nil {
				return fmt.Errorf("failed to create purchase meta data: %w", err)
			}
		} else {
			// Some other error occurred
			return fmt.Errorf("failed to query purchase meta data: %w", err)
		}
	} else {
		// Record exists, update it
		existingPurchaseMetaData.Status = string(model.StatusInCart)
		if err := db.Save(&existingPurchaseMetaData).Error; err != nil {
			return fmt.Errorf("failed to update purchase meta data: %w", err)
		}
	}

	// Check if the delivery status already exists
	// var existingDeliveryStatus model.DeliveryStatus
	// if err := db.Where("cart_item_id = ?", cartItemID).First(&existingDeliveryStatus).Error; err != nil {
	// 	if gorm.ErrRecordNotFound == err {
	// 		// Record doesn't exist, create a new one
	// 		deliveryStatus := &model.DeliveryStatus{
	// 			CartItemID: cartItemID,
	// 			Status:     string(model.StatusOrder),
	// 		}
	// 		if err := db.Create(&deliveryStatus).Error; err != nil {
	// 			return fmt.Errorf("failed to create delivery status: %w", err)
	// 		}
	// 	} else {
	// 		// Some other error occurred
	// 		return fmt.Errorf("failed to query delivery status: %w", err)
	// 	}
	// } else {
	// 	// Record exists, update it
	// 	existingDeliveryStatus.Status = string(model.StatusOrder)
	// 	if err := db.Save(&existingDeliveryStatus).Error; err != nil {
	// 		return fmt.Errorf("failed to update delivery status: %w", err)
	// 	}
	// }

	return nil
}


func GetUserDetails(db *gorm.DB, userUUID uuid.UUID) (userID, userNames, userNumbers string, err error) {
	// Check if the user is a regular user
	var user model.User
	if err := db.Where("id = ?", userUUID).First(&user).Error; err == nil {
		return user.ID.String(), user.Name, user.Number, nil
	}

	// Check if the user is a vendor
	var vendor model.Vendor
	if err := db.Where("id = ?", userUUID).First(&vendor).Error; err == nil {
		return vendor.ID.String(), vendor.Name, vendor.Number, nil
	}

	// Check if the user is a side swapper
	var sideSwapper model.SideSwapper
	if err := db.Where("id = ?", userUUID).First(&sideSwapper).Error; err == nil {
		return sideSwapper.ID.String(), sideSwapper.Name, sideSwapper.Number, nil
	}

	// If no user is found in any of the tables
	return "", "", "", fmt.Errorf("user not found")
}



// handleCartItemGrouping is responsible for handling the grouping of cart items
func  HandleCartItemGrouping(db *gorm.DB,cartItem *model.CartItem) error {
	// Check if the CartItemGroup exists based on the GroupLabel
	var cartItemGroup model.CartItemGroup
	if err := db.Where("group_label = ?", "unGroup").First(&cartItemGroup).Error; err != nil {
		// Group does not exist, create a new CartItemGroup
		cartItemGroup = model.CartItemGroup{
			GroupLabel: "unGroup",
		}
		if err := db.Create(&cartItemGroup).Error; err != nil {
			return fmt.Errorf("failed to create new cart item group: %w", err)
		}
	}

	// Check if the Group already exists for this CartItemID and CartItemGroupID
	var existingGroup model.Group
	if err := db.Where("cart_item_group_id = ? AND cart_item_id = ?", cartItemGroup.ID, cartItem.ID).First(&existingGroup).Error; err == nil {
		// If group exists, increment the QuantityGroup
		existingGroup.QuantityGroup += 1

		// Save the updated group
		if err := db.Save(&existingGroup).Error; err != nil {
			return fmt.Errorf("failed to update existing group: %w", err)
		}
	} else {
		// If group doesn't exist, create a new one
		newGroup := model.Group{
			CartItemID:      cartItem.ID,
			CartItemGroupID: cartItemGroup.ID,
			QuantityGroup:   1,
		}

		// Save the new group in the database
		if err := db.Create(&newGroup).Error; err != nil {
			return fmt.Errorf("failed to create new group: %w", err)
		}
	}

	return nil
}
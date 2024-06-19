package helper

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/mineracail/zokus/graph/model"
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


func CheckDuplicateRecords(db *gorm.DB, model interface{}, inputName, inputEmail, inputPhoneNumber string) error {

	// Check if a member with the same name and phone number already exists
	err := db.Where("name = ? AND phone_number = ?", inputName, inputPhoneNumber).First(&model).Error
	if err == nil {
		return fmt.Errorf("a member with the same name and phone number already exists")
	}
	// else if !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return err
	// }

	// Check if an email already exists
	err = db.Where("email = ?", inputEmail).First(&model).Error
	if err == nil {
		return fmt.Errorf("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err

	}

	// Check if a phone_number already exists
	err = db.Where("phone_number = ?", inputPhoneNumber).First(&model).Error
	if err == nil {
		return fmt.Errorf("phone_number already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
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
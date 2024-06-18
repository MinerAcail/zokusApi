package helper

import (
	"fmt"
	"os"
	"reflect"

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

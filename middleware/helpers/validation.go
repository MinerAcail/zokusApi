package helpers

import (
	"errors"
	"fmt"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

var JwtSecret = []byte("happytime") // Replace this with your own secret key

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func VerifyPassword(hashedPassword string, inputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		return errors.New("email or password is incorrect")
	}
	return nil
}

func GenerateRandomPassword(length int) string {
	// Define a character set for the password
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, length)

	// Generate the password
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}

	return string(password)
}

func GenerateRandomEmail(domain string, username string) string {
	// username := GenerateRandomPassword(6) // Generate a random username
	return fmt.Sprintf("%s@%s", username, domain)
}

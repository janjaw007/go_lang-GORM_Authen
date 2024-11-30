package main

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

func createUser(db *gorm.DB, user *User) error {

	// hashPassword
	hashedPassword, err :=
		bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func loginUser(db *gorm.DB, loginUser *User) (string, error) {
	// get email from user
	user := new(User)
	result := db.Where("email = ?", loginUser.Email).First(user)
	if result.Error != nil {
		return "", result.Error
	}
	// compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))

	if err != nil {
		return "", err
	}
	// pass =  return jwt
	// Create JWT token
	jwtSecretKey := "TestSecret" // should be env
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}
	return t, nil

}

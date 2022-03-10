package utils

import (
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func GenSecretKey() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("Can't generate key", err)
	}
	return key
}

func GetJWT(secKey []byte, userId int) string {
	claims := jwt.StandardClaims{
		ExpiresAt: 15000,
		Id:        strconv.Itoa(userId),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secKey)
	if err != nil {
		fmt.Println("Can't create token due to", err)
	}
	return signedToken
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

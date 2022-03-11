package utils

import (
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type UserClaims struct {
	UserId int
	jwt.StandardClaims
}

func GenSecretKey() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("Can't generate key", err)
	}
	err = os.Setenv("SECRET_KEY", string(key))
	if err != nil {
		fmt.Println("Can't set key to ENV", err)
	}
	return key
}

func GetJWT(secKey []byte, userId int) string {
	claims := UserClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secKey)
	if err != nil {
		fmt.Println("Can't create token due to", err)
	}
	return signedToken
}

func ParseJWT(tokenString string, secKey []byte) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secKey, nil
	})

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		fmt.Println(err)
		return nil, err
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

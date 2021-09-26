package auth

import (
	//"ustore-server/constants"
	"errors"
	"fmt"
	"github.com/google/martian/log"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

func ValidateHeader(bearerHeader string) (interface{}, error) {
	bearerToken := strings.Split(bearerHeader, " ")[1]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error decoding token")
		}
		return []byte("123123123123123"), nil
	})
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	if token.Valid {
		return claims["user"].(string), nil
	}
	return nil, errors.New("invalid token")
}

func GenerateJWT(userEmail string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = userEmail
	claims["exp"] = time.Now().Add(time.Minute * 300).Unix()
	/*
	 Please note that in real world, we need to move "some_secret_key_val_123123" into something like
	 "secret.json" file of Kubernetes etc
	*/
	tokenString, err := token.SignedString([]byte("123123123123123"))
	if err != nil {
		fmt.Println("reached")
		return "", err
	}
	return tokenString, nil
}
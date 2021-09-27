package service

import (
	"database/sql"
	"fmt"
	"ustore/gen/models"
	"ustore/service/auth"
	//	"errors"
	"golang.org/x/crypto/bcrypt"

)

func (c *service) UserLogin(db *sql.DB, email string, password string) (string, error) {
	row := db.QueryRow("SELECT email, password from user where email=?", email)
	userInfo := models.Login{}

	err := row.Scan(&userInfo.Email,
		&userInfo.Password)
	if err != nil {
		return "", err
	}

	//decrypt the hashed-password and compare
	err = bcrypt.CompareHashAndPassword([]byte(*userInfo.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		return "",err
	}
	token, err := auth.GenerateJWT(email)
	if err != nil {
		fmt.Println("error defining token")
		return "",err
	}
	return token, nil

}
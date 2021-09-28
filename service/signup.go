package service

import (
	"database/sql"
	"ustore/gen/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
)


func (c *service) Registration(db *sql.DB, userInfo *models.SignUp) error {
	password := []byte(*userInfo.Password)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	//	err := User.AddUser(db, userInfo)
	row, err := db.Exec(
		"INSERT into user (email, first_name, last_name, middle_name, password, profile_image, username) values (?,?,?,?,?,?,?)",
		userInfo.Email,
		userInfo.FirstName,
		userInfo.LastName,
		userInfo.MiddleName,
		hashedPassword,
		userInfo.ProfileImage,
		userInfo.Username,
	)
	if err != nil {
		return err
	}
	if count, _ := row.RowsAffected(); count != 1 {
		return errors.New("Error inserting row value")
	}

	return nil
}
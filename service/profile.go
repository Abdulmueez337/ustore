package service

import (
	"database/sql"
	"ustore/gen/models"

)

func (c *service) ViewProfile(db *sql.DB, email string) (*models.Profile, error) {
	row := db.QueryRow("SELECT email, username, first_name, middle_name, last_name, profile_image from user where email=?", email)
	userInfo := models.Profile{}

	err := row.Scan(&userInfo.Email,
		&userInfo.Username,
		&userInfo.FirstName,
		&userInfo.MiddleName,
		&userInfo.LastName,
		&userInfo.ProfileImage)
	if err != nil {
		return &userInfo, err
	}

	return &userInfo, nil

}
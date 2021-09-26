package service

import (
	"database/sql"
	"ustore/gen/models"
)

type ServiceInfoHandler interface {
	Registration(db *sql.DB, userInfo *models.SignUp) error
	UserLogin(db *sql.DB, email string, password string) (string, error)
}

type service struct{}

func NewServiceInfoHandler() ServiceInfoHandler{
	return &service{}
}
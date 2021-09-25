package service

import (
	"database/sql"
	"ustore/gen/models"
)

type ServiceInfoHandler interface {
	Registration(db *sql.DB, userInfo *models.SignUp) error
}

type service struct{}

func NewServiceInfoHandler() ServiceInfoHandler{
	return &service{}
}
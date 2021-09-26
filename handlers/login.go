package handlers

import (
	"database/sql"
	"ustore/gen/models"
	"ustore/gen/restapi/operations/login"
	"ustore/service"
	"fmt"
	"github.com/go-openapi/runtime/middleware"

)

type Login struct {
	dbClient            *sql.DB
	serviceInfoHandler service.ServiceInfoHandler
}

func NewLoginHandler(db *sql.DB, serviceInfoHandler service.ServiceInfoHandler) login.LoginHandler {
	return &Login{
		dbClient: db,
		serviceInfoHandler: serviceInfoHandler,
	}
}

func (impl *Login) Handle(params login.LoginParams) middleware.Responder {
	email := *params.Login.Email
	password := *params.Login.Password
	token, err := impl.serviceInfoHandler.UserLogin(impl.dbClient, email, password)
	if err != nil {
		fmt.Println(err.Error())
		return login.NewLoginInternalServerError().WithPayload("Error fetching user details")
	}
	return login.NewLoginOK().WithPayload(&models.LoginSuccess{Success: true, Token: token})
}
package handlers

import (
	"database/sql"
	"ustore/gen/models"
	"ustore/gen/restapi/operations/signup"
	"ustore/service"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"strings"
)

type SignUp struct {
	dbClient            *sql.DB
	serviceInfoHandler service.ServiceInfoHandler
}

func NewSignUpHandler(db *sql.DB, serviceInfoHandler service.ServiceInfoHandler) signup.SignupHandler {
	return &SignUp{
		dbClient:            db,
		serviceInfoHandler : serviceInfoHandler,
	}
}

func (s *SignUp) Handle(params signup.SignupParams) middleware.Responder {
	err := s.serviceInfoHandler.Registration(s.dbClient, params.Signup)
	if err != nil {
		fmt.Println(err.Error())
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "username_UNIQUE") {
				return signup.NewSignupOK().WithPayload(&models.SignUpResponse{Success: false, Message: "Username already registered"})
			}
			if strings.Contains(err.Error(), "email_UNIQUE") {
				return signup.NewSignupOK().WithPayload(&models.SignUpResponse{Success: false, Message: "Email already registered"})
			}
		}
		return signup.NewSignupInternalServerError().WithPayload("Error registering user")
	}
	return signup.NewSignupOK().WithPayload(&models.SignUpResponse{Success: true, Message: "User Registered successfully"})

}
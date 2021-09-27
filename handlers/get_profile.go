package handlers

import (
	"database/sql"
//	"ustore/gen/models"
	"ustore/gen/restapi/operations/user"
	"ustore/service"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"ustore/service/auth"
)

type Profile struct {
	dbClient            *sql.DB
	serviceInfoHandler service.ServiceInfoHandler
}

func NewProfileHandler(db *sql.DB, serviceInfoHandler service.ServiceInfoHandler) user.ProfileHandler{
	return &Profile{
		dbClient: db,
		serviceInfoHandler: serviceInfoHandler,
	}
}

func (p *Profile) Handle(params user.ProfileParams, principal interface{}) middleware.Responder {
	email, err := auth.ValidateHeader(params.HTTPRequest.Header.Get("Authorization"))
	if err != nil {
		return user.NewProfileInternalServerError().WithPayload("error in parsing token")
	}
    userInfo, err := p.serviceInfoHandler.ViewProfile(p.dbClient, email.(string))
	if err != nil {
		fmt.Println(err.Error())
		return user.NewProfileInternalServerError().WithPayload("Error fetching user details")
	}
	return user.NewProfileOK().WithPayload(userInfo)
}

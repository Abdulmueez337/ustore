package handlers

import (
	"database/sql"
//	"ustore/gen/models"
	"ustore/gen/restapi/operations/item"
	"ustore/gen/restapi/operations/user"
	"ustore/service"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"ustore/service/auth"
)

type Item struct {
	dbClient            *sql.DB
	serviceInfoHandler service.ServiceInfoHandler
}

func NewItemHandler(db *sql.DB, serviceInfoHandler service.ServiceInfoHandler) item.ItemsHandler {
	return &Item{
		dbClient: db,
		serviceInfoHandler: serviceInfoHandler,
	}
}

func (i *Item) Handle(params item.ItemsParams, principal interface{}) middleware.Responder {
	_, err := auth.ValidateHeader(params.HTTPRequest.Header.Get("Authorization"))
	if err != nil {
		return item.NewItemsInternalServerError().WithPayload("error in parsing items")
	}
	products, err := i.serviceInfoHandler.ViewItems(i.dbClient)
	if err != nil {
		fmt.Println(err.Error())
		return user.NewProfileInternalServerError().WithPayload("Error fetching items details")
	}
	return item.NewItemsOK().WithPayload(products)
}

package handlers

import (
	"database/sql"
	"ustore/gen/models"
	"ustore/gen/restapi/operations/item"
	"ustore/service"
	"github.com/go-openapi/runtime/middleware"
	"ustore/service/auth"
)

type Subscription struct {
	dbClient            *sql.DB
	serviceInfoHandler service.ServiceInfoHandler
}

func NewSubscriptionHandler(db *sql.DB, serviceInfoHandler service.ServiceInfoHandler) item.SubscribeHandler {
	return &Subscription{
		dbClient:            db,
		serviceInfoHandler : serviceInfoHandler,
	}
}

func (s *Subscription) Handle(params item.SubscribeParams, principal interface {}) middleware.Responder {
	email, err := auth.ValidateHeader(params.HTTPRequest.Header.Get("Authorization"))
	if err != nil {
		return item.NewSubscribeInternalServerError().WithPayload("error in parsing token")
	}
	err = s.serviceInfoHandler.SubscribeItem(s.dbClient, email.(string), params.Subscribe)
	if err != nil {
		return item.NewSubscribeOK().WithPayload(&models.SubscriptionResponse{Success: false, Message: "Item could not be subscribed"})
	}

	return item.NewSubscribeOK().WithPayload(&models.SubscriptionResponse{Success: true, Message: "Item Subscribed Successfully"})

}
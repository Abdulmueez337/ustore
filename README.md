# Get Items - EntryPoint

The fourth entrypoint is ViewItems.
To generate code for ViewItems API, swagger.yml file is updated accordingly for path and model definitions.
<b>swagger.yml</b>
```
 //path
 /items:
    get:
      description: "To show available items"
      operationId: "items"
      tags:
        - "item"
      security:
        - Bearer: [ ]
      responses:
        200:
          description: "Success response when items are shown"
          schema:
            $ref: "#/definitions/Products"
        400:
          description: Bad Request
        404:
          description: items not found
        500:
          schema:
            type: string
          description: Server error
          
// model definition
  Products:
    type: array
    items:
      $ref: "#/definitions/Product"
  Product:
    type: object
    properties:
      item_name:
        type: string
      item_details:
        type: string
      monthly_price:
        type: number
      yearly_price:
        type: number
      available_items:
        type: integer
```

Now we again generate the code in gen directory for the ViewItems entrypoint with following command.
```
swagger generate server -t gen -f ./swagger.yml --default-scheme http --exclude-main

// to install dependencies
go get -u -f ./gen/...
```

As we are directly hitting database from service layer, Thus item.go file is created in service directory and implemented ViewItems function. Also add the ViewItems function in the service interface in <b>ustore/service/service.go</b>
<b>ustore/service/items.go</b>
```
package service

import (
	"database/sql"
	"ustore/gen/models"

)

func (c *service) ViewItems(db *sql.DB) (models.Products, error) {
	rows, _ := db.Query("SELECT item_name, item_details, monthly_price, yearly_price, available_items from item")
	products:= models.Products{}

	for rows.Next() {
		item := models.Product{}
		err := rows.Scan(
			&item.ItemName,
			&item.ItemDetails,
			&item.MonthlyPrice,
			&item.YearlyPrice,
			&item.AvailableItems,
			)
		if err != nil {
			return products, err
		}
		products = append(products, &item)
	}
	return products, nil
}
```

The items handler hits the ViewItems function in service layer as following.
```
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
```
Finally, the handler is called in configureAPI function <b>ustore/gen/restapi/operation/configure_ustore.go</b>
```
api.ItemItemsHandler = handlers.NewItemHandler(db, serviceInfoHandle)
```

#### Testing

![](https://i.imgur.com/6ANVilB.png)




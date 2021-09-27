# ViewProfile - EntryPoint

The third entrypoint is ViewUserProfile.
To generate code for ViewProfile API, swagger.yml file is updated accordingly for path and model definitions.
<b>swagger.yml</b>
```
 //path
 /user/profile:
    get:
      description: "To show user details"
      operationId: "profile"
      tags:
        - "User"
      security:
        - Bearer: []
      responses:
        200:
          description: "Success response when item is added successfully"
          schema:
            $ref: "#/definitions/Profile"
        400:
          description: Bad Request
        404:
          description: User not found
        500:
          schema:
            type: string
          description: Server error
          
// model definition
  Profile:
    type: object
    properties:
      first_name:
        type: string
      middle_name:
        type: string
      last_name:
        type: string
      email:
        type: string
      username:
        type: string
      profile_image:
        type: string
```

Now we again generate the code in gen directory for the ViewProfile entrypoint with following command.
```
swagger generate server -t gen -f ./swagger.yml --default-scheme http --exclude-main

// to install dependencies
go get -u -f ./gen/...
```

As we are directly hitting database from service layer, Thus profile.go file is created in service directory and implemented ViewProfile function. Also add the ViewProfile function in the service interface in <b>ustore/service/service.go</b>
<b>ustore/service/profile.go</b>
```
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
```

The profile handler hits the ViewProfile function in service layer as following.
```
package handlers

import (
	"database/sql"
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

```
Finally, the handler is called in configureAPI function <b>ustore/gen/restapi/operation/configure_ustore.go</b>
```
api.UserProfileHandler = handlers.NewProfileHandler(db, serviceInfoHandle)
```

#### Testing

![](https://i.imgur.com/HtARYwC.png)



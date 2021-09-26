# Login - EntryPoint

We continue the same signup entrypoint and add login api for the ustore microservice.
First we generate code for our login entrypoint as we did for signup api.
<b>swagger.yml</b>
```
// defining JWT for later use
securityDefinitions:
  Bearer:
    type: apiKey
    name: Authorization
    in: header
    
// login path
 /login:
    post:
      description: 'Returns token for authorized User'
      operationId: "login"
      tags:
        - "login"
      parameters:
        - in: 'body'
          name: 'login'
          required: true
          description: "login model"
          schema:
            $ref: "#/definitions/Login"
      responses:
        200:
          description: Successful login
          schema:
            $ref: '#/definitions/LoginSuccess'
        400:
          description: Bad Request
        404:
          schema:
            type: string
          description: User not found
        500:
          schema:
            type: string
          description: Server error

// Login and LoginSuccess model definitions
  LoginInfo:
    type: object
    required: [email,password]
    properties:
      email:
        type: string
      password:
        type: string
  LoginSuccess:
    type: object
    properties:
      success:
        type: boolean
      token:
        type: string
```

Now we again generate the code in gen directory for our login entrypoint with following command.
```
swagger generate server -t gen -f ./swagger.yml --default-scheme http --exclude-main

// to install dependencies
go get -u -f ./gen/...
```
Before implementing the login in service layer, authentication JWT is required to implement in the service layer.
<b>ustore/service/auth/autherization.go</b>
```
package auth

import (
	//"ustore-server/constants"
	"errors"
	"fmt"
	"github.com/google/martian/log"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

func ValidateHeader(bearerHeader string) (interface{}, error) {
	bearerToken := strings.Split(bearerHeader, " ")[1]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error decoding token")
		}
		return []byte("123123123123123"), nil
	})
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	if token.Valid {
		return claims["user"].(string), nil
	}
	return nil, errors.New("invalid token")
}

func GenerateJWT(userEmail string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = userEmail
	claims["exp"] = time.Now().Add(time.Minute * 300).Unix()
	/*
	 Please note that in real world, we need to move "some_secret_key_val_123123" into something like
	 "secret.json" file of Kubernetes etc
	*/
	tokenString, err := token.SignedString([]byte("123123123123123"))
	if err != nil {
		fmt.Println("reached")
		return "", err
	}
	return tokenString, nil
}
```

As we are directly hitting database from service layer, Thus login.go file is created in service directory and implement the login function. Also include the UserLogin function in the service interface in <b>ustore/service/service.go</b>
<b>ustore/service/login.go</b>
```
package service

import (
	"database/sql"
	"fmt"
	"ustore/gen/models"
	"ustore/service/auth"
	//	"errors"
	"golang.org/x/crypto/bcrypt"

)

func (c *service) UserLogin(db *sql.DB, email string, password string) (string, error) {
	row := db.QueryRow("SELECT email, password from user where email=?", email)
	userInfo := models.Login{}

	err := row.Scan(&userInfo.Email,
		&userInfo.Password)
	if err != nil {
		return "", err
	}

	//decrypt the hashed-password and compare
	err = bcrypt.CompareHashAndPassword([]byte(*userInfo.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		return "",err
	}
	token, err := auth.GenerateJWT(email)
	if err != nil {
		fmt.Println("error defining token")
		return "",err
	}
	return token, nil

}
```

The login-handler hits the UserLogin function in service layer as following.
```
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
```
Finally, the handler is called in configureAPI function <b>ustore/gen/restapi/operation/configure_ustore.go</b>
```

    api.LoginLoginHandler = handlers.NewLoginHandler(db, serviceInfoHandle)
```

#### Testing


![](https://i.imgur.com/VsQbCAi.png)

![](https://i.imgur.com/cIRsPch.png)



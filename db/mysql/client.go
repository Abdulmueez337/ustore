package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type client struct {
	// this place can be used to initialize the auth client which can be used to talk to other micro-services
}

func NewClient() client {
	return client{}
}

func (b client) BuildSqlClient() *sql.DB {
	// sensitive info can be stored in "secrets.json" of GKE
	db, err := sql.Open("mysql", "simsim:MYpassword100@/ustore?parseTime=True")
	if err != nil {
		log.Fatal("error connecting DB : ", err.Error())
	}
	return db
}
package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"ustore/config"
)

type client struct {
	// this place can be used to initialize the auth client which can be used to talk to other micro-services
}

func NewClient() client {
	return client{}
}

func (b client) BuildSqlClient() *sql.DB {
	// database connection configuration
	db, err := sql.Open("mysql", config.UserName+":"+config.Password+"@/"+config.DbSchema+"?parseTime=True" )
	if err != nil {
		log.Fatal("error connecting DB : ", err.Error())
	}
	return db
}
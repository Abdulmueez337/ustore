package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"ustore/gen/models"

)

func (c *service) SubscribeItem(db *sql.DB, email string, subInfo *models.Subscribe) error {
	row := db.QueryRow("SELECT user_id from user where email=?", email)
	var user_id *int
	err := row.Scan(&user_id)
	if err != nil {
		fmt.Println("userid issue")
		return err
	}
	//UTC time : Y-M-D h-m-s
	startTime := time.Now().UTC()
    endTime := time.Now().AddDate(0,0,30).UTC()
	r, err := db.Exec(
		"INSERT into subscription (start_time, end_time,subs_price, status,user_id,item_id ) values (?,?,?,?,?,?)",
		startTime,
		endTime,
		subInfo.SubsPrice,
		subInfo.Status,
		*user_id,
		subInfo.ItemID,
	)


	if err != nil {
		fmt.Println("errror in insertion", *user_id, err)
		return err
	}
	if count, _ := r.RowsAffected(); count != 1 {
		return errors.New("Error inserting row value")
	}
	return nil

}
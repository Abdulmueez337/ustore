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
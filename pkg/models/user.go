package models

type User struct {
	Id              int    `db:"user_id"`
	Name            string `db:"user_name"`
	Number          string `db:"user_number"`
	City            string `db:"user_city"`
	DeliveryAddress string `db:"delivery_adress"`
}

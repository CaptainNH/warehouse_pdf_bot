package models

type Courier struct {
	Id          int    `db:"id"`
	Name        string `db:"name"`
	PhoneNumber string `db:"number"`
}

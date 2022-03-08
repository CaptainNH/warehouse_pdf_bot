package models

type ShopingCart struct {
	OrderId        int    `db:"order_id"`
	ProductId      int    `db:"product_id"`
	Price          int    `db:"price"`
	UnitPrice      int    `db:"unit_price"`
	DeliveryFormat string `db:"delivery_format"`
	Quantity       int    `db:"quantity"`
}

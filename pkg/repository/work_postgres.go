package repository

import (
	"fmt"
	"gopls-workspace/pkg/models"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetOrders(db *sqlx.DB, date string) []int {
	numbers := []int{}
	//err := db.Select(&numbers, fmt.Sprintf("SELECT id FROM orders WHERE order_date LIKE '%s' AND order_status != 'denied'", "20.02.2022%"))
	err := db.Select(&numbers, fmt.Sprintf("SELECT order_id FROM couriers_orders WHERE date LIKE '%s'", date))
	if err != nil {
		log.Fatalf("error select order id: %s", err.Error())
	}
	return numbers
}

func GetShopingCart(db *sqlx.DB, orderId int) []models.ShopingCart {
	sc := []models.ShopingCart{}
	err := db.Select(&sc, fmt.Sprintf("SELECT * FROM shoping_cart WHERE order_id = %d", orderId))
	if err != nil {
		log.Fatalf("error select shoping cart: %s", err.Error())
	}
	return sc
}

func GetProductTitle(db *sqlx.DB, productId int) string {
	var productTitle string
	err := db.Get(&productTitle, fmt.Sprintf("SELECT title FROM products WHERE id = %d", productId))
	if err != nil {
		log.Fatalf("error get product title: %s", err.Error())
	}
	return productTitle
}

func GetProductsList(db *sqlx.DB, orderId int) []string {
	shopingCart := GetShopingCart(db, orderId)
	productsList := []string{}
	for i, sc := range shopingCart {
		productTitle := GetProductTitle(db, sc.ProductId)
		productsList = append(productsList, fmt.Sprintf("%d. %s x%d %s %.2fруб.",
			i+1, productTitle, sc.Quantity, sc.DeliveryFormat, float64(sc.Price)/100.0))
	}
	return productsList
}

func GetFullPrice(db *sqlx.DB, orderId int) float64 {
	shopingCart := GetShopingCart(db, orderId)
	var sum float64
	for _, sc := range shopingCart {
		sum += (float64(sc.Price) / 100.0)
	}
	return sum
}

func GetUserData(db *sqlx.DB, orderId int) models.User {
	var user models.User
	err := db.Get(&user,
		fmt.Sprintf("SELECT user_id, user_name, user_number, user_city, delivery_adress FROM orders WHERE id = %d", orderId))
	if err != nil {
		log.Fatalf("error get user info: %s", err.Error())
	}
	return user
}

func GetCourierId(db *sqlx.DB, orderId int) int {
	var courierId int
	err := db.Get(&courierId, fmt.Sprintf("SELECT courier_id FROM couriers_orders WHERE order_id = %d", orderId))
	if err != nil {
		log.Fatalf("error get courier id: %s", err.Error())
	}
	return courierId
}

func GetCourierData(db *sqlx.DB, orderId int) models.Courier {
	courierId := GetCourierId(db, orderId)
	var courier models.Courier
	err := db.Get(&courier, fmt.Sprintf("SELECT id, name, number FROM couriers WHERE id = %d", courierId))
	if err != nil {
		log.Fatalf("error get courier info: %s", err.Error())
	}
	return courier
}

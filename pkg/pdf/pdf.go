package pdf

import (
	"fmt"
	"gopls-workspace/pkg/repository"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func InitDB() *sqlx.DB {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	return db
}

func CreateFile() (string, bool) {
	db := InitDB()
	date := time.Now()
	filePath := fmt.Sprintf("files/Заказы %s.pdf", date.Format("02.01.2006"))
	ordersId := repository.GetOrders(db, date.Format("02.01.2006")+"%")
	if len(ordersId) == 0 {
		return "", false
	}
	log.Print("Creating file...")
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddFont("Helvetica", "", "configs/helvetica_1251.json")
	pdf.SetFont("Helvetica", "", 16)
	tr := pdf.UnicodeTranslatorFromDescriptor("configs/cp1251")
	for _, id := range ordersId {
		productsList := repository.GetProductsList(db, id)
		fullPrice := repository.GetFullPrice(db, id)
		userData := repository.GetUserData(db, id)
		courierData := repository.GetCourierData(db, id)
		pdf.AddPage()
		pdf.MultiCell(200, 10, tr(fmt.Sprintf("Заказ № %d\nСписок продуктов:\n", id)), "", "", false)
		pdf.MultiCell(200, 10, tr(strings.Join(productsList, "\n")), "", "", false)
		pdf.MultiCell(200, 10, tr(fmt.Sprintf("Общая сумма заказа: %.2f\nДанные о заказчике:\n", fullPrice)), "", "", false)
		pdf.MultiCell(200, 10, tr(fmt.Sprintf("Имя заказчика: %s\nНомер заказчика: %s\nФилиал: %s\nАдрес доставки: %s",
			userData.Name, userData.Number, userData.City, userData.DeliveryAddress)), "", "", false)
		pdf.MultiCell(200, 10, tr(fmt.Sprintf("Данные о курьере:\nИмя курьера: %s\nНомер курьера: %s",
			courierData.Name, courierData.PhoneNumber)), "", "", false)
	}
	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		log.Fatalf("failed to create file: %s", err.Error())
	}
	return filePath, true
}

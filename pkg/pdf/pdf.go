package pdf

import (
	"fmt"
	"gopls-workspace/pkg/repository"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	gofpdf "github.com/jung-kurt/gofpdf"
	"github.com/spf13/viper"
)

func InitDB() *sqlx.DB {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.password"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	return db
}

func CreateFile() {
	log.Print("Creating file...")
	id := 1
	productsList := "1. Помидоры x1 200.0р\n2. Огурцы x1 250.0р"
	sum := "3142р"
	userData := "jfvxksd"
	courierData := "fweff"
	orderInfo := fmt.Sprintf("Заказ № %d\nСписок продуктов:\n%s\nОбщая сумма заказа: %s\nДанные о заказчике:\n%s\nДанные о курьере:\n%s",
		id, productsList, sum, userData, courierData)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 16)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	pdf.MultiCell(70, 10, tr(orderInfo), "", "", false)
	err := pdf.OutputFileAndClose("test.pdf")
	if err != nil {
		log.Fatalf("failed to create file: %s", err.Error())
	}
}

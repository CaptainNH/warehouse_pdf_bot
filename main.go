package main

import (
	"gopls-workspace/pkg/pdf"
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	// глобальная переменная в которой хранится токен
	telegramBotToken string
	// глобальная переменная в которой хранится id чата
	chatID int64
)

func initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	telegramBotToken = os.Getenv("TOKEN")
	id, _ := strconv.Atoi(os.Getenv("CHATID"))
	chatID = int64(id)

	if telegramBotToken == "" || chatID == 0 {
		log.Print("-telegrambottoken and chatID is required")
		os.Exit(1)
	}
}

func sendFile(bot *tgbotapi.BotAPI) {
	pdf.CreateFile()
	log.Print("done!")
	//msg := tgbotapi.NewMessage(chatID, "Заходит скелет в бар...")
	msg := tgbotapi.NewDocument(chatID, tgbotapi.FilePath("test.pdf"))
	bot.Send(msg)
}

func monitor(bot *tgbotapi.BotAPI) {
	for {
		sendFile(bot)
		time.Sleep(time.Minute * 5)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	initialize()
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	go monitor(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// универсальный ответ на любое сообщение
		reply := "Введите команду /request"
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Command() == "request" {
			sendFile(bot)
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}

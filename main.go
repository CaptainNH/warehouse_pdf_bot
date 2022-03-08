package main

import (
	"fmt"
	"gopls-workspace/pkg/pdf"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var (
	telegramBotToken string
	chatID           int64
	fileSendTime     string
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
	filePath, ok := pdf.CreateFile()
	if !ok {
		log.Print("orders is empty")
		msg := tgbotapi.NewMessage(chatID, "Нет новых заказов")
		bot.Send(msg)
		return
	}
	log.Print("done!")
	msg := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(filePath))
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatalf("error sending file: %s", err.Error())
	}
}

func monitor(bot *tgbotapi.BotAPI) {
	for {
		if fileSendTime == time.Now().Format("15:04") {
			sendFile(bot)
			time.Sleep(time.Minute * 5)
		}
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	fileSendTime = "17:04"
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
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Установите время отправки файла командой /settime\nЧтобы получить файл прямо сейчас используйте команду /request")
			bot.Send(msg)
		case "request":
			sendFile(bot)
		case "settime":
			var reply string
			cmd := strings.Split(update.Message.Text, " ")
			if len(cmd) > 1 {
				t, err := time.Parse("15:04", strings.Split(update.Message.Text, " ")[1])
				if err != nil {
					reply = "Введите корректное время в формате \"HH:mm\""
				}
				fileSendTime = t.Format("15:04")
				reply = fmt.Sprintf("Файл будет отправляться в %s", fileSendTime)
			} else {
				reply = "Введите время в формате \"HH:mm\""
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			bot.Send(msg)
		}
	}
}

package main

import (
	"log"
	//"strconv"
	"time"
	//ntp "github.com/beevik/ntp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	// глобальная переменная в которой храним токен
	telegramBotToken string = "5142340650:AAFDMa-fUuLHTuV1NvqtmwwFOBYj4HPMS_E"
	chatID           int64  = -775300667
)

// func init() {
// 	// принимаем на входе флаг -telegrambottoken
// 	flag.StringVar(&telegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
// 	flag.Parse()

// 	// без него не запускаемся
// 	if telegramBotToken == "" {
// 		log.Print("-telegrambottoken is required")
// 		os.Exit(1)
// 	}
// }

func sendFile(bot *tgbotapi.BotAPI) {
	for {
		msg := tgbotapi.NewMessage(chatID, "Заходит скелет в бар...")
		bot.Send(msg)
		time.Sleep(time.Minute * 5)
	}
}

func main() {
	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	go sendFile(bot)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	for update := range updates {
		// универсальный ответ на любое сообщение
		reply := "Не знаю что сказать"
		if update.Message == nil {
			continue
		}

		// логируем от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// свитч на обработку комманд
		// комманда - сообщение, начинающееся с "/"
		switch update.Message.Command() {
		case "start":
			reply = "Привет. Я телеграм-бот"
		case "hello":
			reply = "world"
		case "request":
			reply = time.Now().Format("15:04:05")
		}

		// создаем ответное сообщение
		msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		msg := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FilePath("C:\\Users\\777\\Desktop\\Трушин.pdf"))
		msg.ReplyToMessageID = update.Message.MessageID
		// отправляем

		bot.Send(msg)
		bot.Send(msg1)
	}
}

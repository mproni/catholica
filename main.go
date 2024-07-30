package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	// включаем сообщения отладки
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// конфиг обновления
	// tgbotapi.NewUpdate(0) - означает с какого необработанного сообщения
	// начать обработку, например, если мы отключили бота на полчасика, а
	// нам пришли за это время сообщения, то Offset в данном случае нам
	// говорит: "начни читать с первого необработанного сообщения"
	//
	// чтобы пропустить все необработанные сообщения, то нам всё равно
	// нужно начать с 0 и все их прочитать, но пропустить обработку
	updateConfig := tgbotapi.NewUpdate(0)
	// таймаут работает следующим образом. мы отправляем запрос на
	// открытие соединения с Телеграмом и в течение этого времени
	// Телеграм чекает у себя есть ли новые сообщения, если такое есть,
	// то он отправляет его нам, а если за это время ничего не произошло,
	// то мы повторно отправляем запрос на открытие нового соединения,
	// чтобы и дальше можно было получать сообщения.
	updateConfig.Timeout = 60
	// получаем канал обновлений, привязывая его к нашему боту,
	// куда будут поступать сообщения с настройками из updateConfig,
	// в данном случае Offset = 0, Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	// как только поступило сообщение, то присваиваем
	// его переменной update и обрабатываем.
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// если отправлять в ответ update.Message.Text, то может возникнуть
		// такая ситуация, что сообщение будет пустым (если нам пришла,
		// например, картинка), это будет Bad Request и мы вылетим с panic.
		msg := tgbotapi.NewMessage(update.Message.From.ID, "Сообщение получено!")

		_, err := bot.Send(msg)
		if err != nil {
			panic(err)
		}
	}
}

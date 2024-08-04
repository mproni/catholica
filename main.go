package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Структуры для обработки данных Telegram
type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

type SendMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// Функция для получения обновлений
func getUpdates(token string, offset int) ([]Update, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d", token, offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Result, nil
}

// Функция для отправки сообщения
func sendMessage(token string, chatID int64, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	message := SendMessage{
		ChatID: chatID,
		Text:   text,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
	}
	var lastUpdateID int

	for {
		updates, err := getUpdates(token, lastUpdateID+1)
		if err != nil {
			log.Println("Error getting updates:", err)
			continue
		}

		for _, update := range updates {
			if update.UpdateID > lastUpdateID {
				lastUpdateID = update.UpdateID
			}

			if update.Message.Text != "" {
				log.Printf("[%s] %s", update.Message.From.Username, update.Message.Text)
				err := sendMessage(token, update.Message.Chat.ID, "Hello, "+update.Message.From.FirstName+"!")
				if err != nil {
					log.Println("Error sending message:", err)
				}
			}
		}
	}
}

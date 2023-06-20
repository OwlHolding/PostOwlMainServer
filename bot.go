package main

import (
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateHandler func(int64, string, string)

var BotAPI *tgbotapi.BotAPI
var BotUpdateHandler UpdateHandler

func InitBot(config ServerConfig, handler UpdateHandler) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatal(err)
	}

	webhook, _ := tgbotapi.NewWebhookWithCert(config.Url+":"+config.Port+"/"+config.Token+"/",
		tgbotapi.FilePath(config.CertFile))
	webhook.MaxConnections = config.MaxBotConns

	_, err = bot.Request(webhook)
	if err != nil {
		log.Fatal(err)
	}

	BotAPI = bot
	BotUpdateHandler = handler

	log.Println("Bot: inited")
}

func ProcessRequest(PostBody []byte) {
	var update tgbotapi.Update
	err := json.Unmarshal(PostBody, &update)
	if err != nil {
		log.Fatal(err)
	}

	go BotUpdateHandler(update.Message.From.ID,
		update.Message.Text, update.Message.Chat.UserName)
}

func SendMessage(chatID int64, text string) {
	message := tgbotapi.NewMessage(chatID, text)
	message.ParseMode = "HTML"
	message.DisableWebPagePreview = true
	_, err := BotAPI.Send(message)
	if err != nil {
		log.Fatal(fmt.Errorf("botapi error: %s", err.Error()))
	}
}

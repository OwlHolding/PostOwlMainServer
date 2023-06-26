package main

import (
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateHandler func(int64, string, string)
type CallbackHandler func(int64, int, string, string)

var BotAPI *tgbotapi.BotAPI
var BotUpdateHandler UpdateHandler
var BotCallbackHandler CallbackHandler

func InitBot(config ServerConfig, handler UpdateHandler, callbackhandler CallbackHandler) {
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
	BotCallbackHandler = callbackhandler

	log.Println("Bot: inited")
}

func ProcessRequest(PostBody []byte) {
	var update tgbotapi.Update
	err := json.Unmarshal(PostBody, &update)
	if err != nil {
		log.Fatal(err)
	}
	if update.CallbackQuery != nil {
		BotCallbackHandler(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID,
			update.CallbackQuery.Data, update.CallbackQuery.Message.Text)
	} else {
		go BotUpdateHandler(update.Message.From.ID,
			update.Message.Text, update.Message.Chat.UserName)
	}
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

func SendMessageWithInlineKeyboard(chatID int64, text string, channel string) {
	message := tgbotapi.NewMessage(chatID, text)
	keyboard := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("👍", "1"+channel),
				tgbotapi.NewInlineKeyboardButtonData("👎", "0"+channel))},
	}
	message.ReplyMarkup = keyboard
	message.ParseMode = "HTML"
	message.DisableWebPagePreview = true
	_, err := BotAPI.Send(message)
	if err != nil {
		panic(fmt.Errorf("botapi error: %s", err.Error()))
	}
}

func DisableInlineKeyboard(chatID int64, messageID int) {
	message := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID,
		tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{}))
	_, err := BotAPI.Send(message)
	if err != nil {
		panic(fmt.Errorf("botapi error: %s", err.Error()))
	}
}

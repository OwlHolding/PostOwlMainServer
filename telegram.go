package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Post struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

var TelegramWorkers []string

func InitTelegram(config ServerConfig) {
	TelegramWorkers = config.TelegramWorkers
}

func ProcessPost(PostBody []byte) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
			SendToAdmins(fmt.Sprintf("PostError: %s", err))
		}
	}()

	var post Post
	err := json.Unmarshal(PostBody, &post)
	if err != nil {
		log.Println(err)
	}

	users := DataBaseGetUsers(post.Channel)
	if len(users) > 0 {
		users = ApiPredict(post.Channel, post.Text, users)

		for _, user := range users {
			post.Text += "\n@" + post.Channel
			SendMessageWithInlineKeyboard(user, post.Text, post.Channel)
		}
	}
}

func TelegramAddChannel(channel string) {
	var url string
	for _, ip := range TelegramWorkers {
		url = ip + "add-channel/" + channel + "/"
		client := http.Client{Timeout: time.Minute}
		resp, _ := client.Post(url, "", nil)
		if resp.StatusCode == 201 || resp.StatusCode == 200 {
			return
		}
	}
	panic(fmt.Sprintf("Can`t add channel %s to any telegram worker", channel))
}

func TelegramDelChannel(channel string) {
	var url string
	for _, ip := range TelegramWorkers {
		url = ip + "del-channel/" + channel + "/"
		client := http.Client{Timeout: time.Minute}
		resp, _ := client.Post(url, "", nil)
		if resp.StatusCode == 205 || resp.StatusCode == 200 {
			return
		}
	}
	panic(fmt.Sprintf("Can`t delete channel %s from any telegram worker", channel))
}

func ChannelIsExist(channel string) bool {
	url := "https://rsshub.app/telegram/channel/" + channel

	client := http.Client{Timeout: 30 * time.Second}
	resp, _ := client.Get(url)
	if resp.StatusCode/10 == 20 {
		return true
	}
	return false
}

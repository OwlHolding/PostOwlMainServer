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
	var post Post
	err := json.Unmarshal(PostBody, &post)
	if err != nil {
		log.Fatal(err)
	}

	users := DataBaseGetUsers(post.Channel)
	users = ApiPredict(post.Channel, post.Text, users)

	for _, user := range users {
		post.Text += "\n@" + post.Channel
		SendMessage(user, post.Text)
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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var MlServers []string

func InitApi(config ServerConfig) {
	servers := config.MlServers
	MlServers = servers
	log.Println("API: inited")
}

func GetMLUrl() string {
	return MlServers[0]
}

func ApiAddUser(userID int64) {
	url := GetMLUrl() + "add-user/" + fmt.Sprint(userID) + "/"

	client := http.Client{Timeout: time.Minute}

	resp, err := client.Post(url, "", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 208 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}
	log.Printf("API: Add user %d - %d \n", userID, resp.StatusCode)
}

func ApiDelUser(userID int64) {
	url := GetMLUrl() + "del-user/" + fmt.Sprint(userID) + "/"
	client := http.Client{Timeout: time.Minute}

	resp, err := client.Post(url, "", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 205 && resp.StatusCode != 208 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}

	log.Printf("API: Del user %d - %d \n", userID, resp.StatusCode)
}

func ApiAddChannel(userID int64, channel string) {
	url := GetMLUrl() + "add-channel/" + fmt.Sprint(userID) + "/" + channel + "/"
	client := http.Client{Timeout: time.Minute}

	resp, err := client.Post(url, "", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 208 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}
	log.Printf("API: Add channel %s to user %d - %d \n", channel, userID, resp.StatusCode)
}

func ApiDelChannel(userID int64, channel string) {
	url := GetMLUrl() + "del-channel/" + fmt.Sprint(userID) + "/" + channel + "/"
	client := http.Client{Timeout: time.Minute}

	resp, err := client.Post(url, "", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 205 && resp.StatusCode != 208 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}
	log.Printf("API: Del channel %s from user %d - %d \n", channel, userID, resp.StatusCode)
}

type PredictRequest struct {
	Post    string  `json:"post"`
	Channel string  `json:"channel"`
	Users   []int64 `json:"users"`
}

type PredictResponse struct {
	Users []int64 `json:"users"`
}

func ApiPredict(channel string, text string, users []int64) []int64 {
	url := GetMLUrl() + "predict/"

	request := PredictRequest{Channel: channel, Post: text, Users: users}

	json_byte, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	client := http.Client{Timeout: 5 * time.Minute}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(json_byte))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		return make([]int64, 0)
	}

	byte_body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("APIPredictError: %s", err.Error()))
	}

	var response PredictResponse
	err = json.Unmarshal(byte_body, &response)
	if err != nil {
		panic(err)
	}

	log.Printf("API: Predict for channel %s - %d \n", channel, resp.StatusCode)
	return response.Users
}

type TrainRequest struct {
	Text  string `json:"text"`
	Label int16  `json:"label"`
}

func ApiTrain(userID int64, channel string, text string, label int16) {
	url := GetMLUrl() + "train/" + fmt.Sprint(userID) + "/" + channel + "/"

	request := TrainRequest{Text: text, Label: label}
	json_byte, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	client := http.Client{Timeout: time.Minute}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(json_byte))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	log.Printf("API: Train for channel %s and user %d - %d \n", channel, userID, resp.StatusCode)
}

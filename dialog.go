package main

import (
	"fmt"
	"log"
	"strings"
)

const (
	StateIdle           = 0
	StateWaitAddChannel = 1
	StateWaitDelChannel = 2
)

var WhiteList []int64
var BanList []int64
var AdminChatIDs []int64
var AccessKeys []string

var MaxChannels int

func InitDialog(config ServerConfig) {
	WhiteList = config.WhiteList
	BanList = config.BanList
	AdminChatIDs = config.AdminChatIDs
	AccessKeys = config.AccessKeys
	MaxChannels = config.ChansPerUser

	log.Println("Dialog: inited")
}

func CheckBan(chatID int64) bool {
	if len(WhiteList) > 0 {
		ban := true
		for _, id := range WhiteList {
			if id == chatID {
				ban = false
				break
			}
		}
		if ban {
			return true
		}
	} else {
		for _, id := range BanList {
			if id == chatID {
				return true
			}
		}
	}
	return false
}

func SendToAdmins(message string) {
	if len(AdminChatIDs) != 0 {
		for _, AdminChatID := range AdminChatIDs {
			if AdminChatID != 0 {
				SendMessage(AdminChatID, message)
			}
		}
	}
}

func ProcessDialog(userID int64, text string, username string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
			SendMessage(userID, MessageError)
			SendToAdmins(fmt.Sprintf(`Error: "%s"; Username: @%s`, err, username))
		}
	}()

	if CheckBan(userID) {
		SendMessage(userID, MessageBanned)
		return
	}

	if !DataBaseIsUserExist(userID) {
		if len(AccessKeys) != 0 {
			allow := false
			for _, key := range AccessKeys {
				if key == text {
					allow = true
					break
				}
			}
			if allow {
				text = "/start"
			} else {
				SendMessage(userID, MessageNotAllowed)
				return
			}

			DataBaseAddUser(userID)
			go ApiAddUser(userID)
		}
	}

	if text == "/start" {
		SetState(userID, StateIdle)
		SendMessage(userID, MessageHello)

	} else if text == "/addchannel" {
		SetState(userID, StateWaitAddChannel)
		SendMessage(userID, MessageAddChannel)

	} else if text == "/delchannel" {
		SetState(userID, StateWaitDelChannel)
		SendMessage(userID, MessageDelChannel)

	} else if text == "/info" {
		channelslist := DataBaseInfo(userID)
		channels := ""
		for _, channel := range channelslist {
			channels += "\t<code>" + channel + "</code>\n"
		}
		SendMessage(userID, fmt.Sprintf(MessageInfo, channels))

	} else if text == "/cancel" {
		SetState(userID, StateIdle)
		SendMessage(userID, MessageCancel)

	} else {

		state := GetState(userID)

		if state == StateWaitAddChannel {

			if len(DataBaseInfo(userID)) >= MaxChannels {
				SendMessage(userID, fmt.Sprintf(MessageChannelOverflow, MaxChannels))
			}
			text = strings.ReplaceAll(text, "https://t.me/", "")
			text = strings.ReplaceAll(text, "@", "")

			if ChannelIsExist(text) {
				if !DataBaseAddChannel(userID, text) {
					SendMessage(userID, MessageChannelAlreadyAdded)
				} else {
					go ApiAddChannel(userID, text)
					if len(DataBaseGetUsers(text)) == 0 {
						go TelegramAddChannel(text)
					}
					SendMessage(userID, fmt.Sprintf(MessageAddChannelOK, text))
				}
			} else {
				SendMessage(userID, MessageChannelNotExists)
			}

			SetState(userID, StateIdle)

		} else if state == StateWaitDelChannel {

			text = strings.ReplaceAll(text, "https://t.me/", "")
			text = strings.ReplaceAll(text, "@", "")

			if !DataBaseDelChannel(userID, text) {
				SendMessage(userID, MessageChannelNotListed)
			} else {
				go ApiDelChannel(userID, text)
				if len(DataBaseGetUsers(text)) == 0 {
					go TelegramDelChannel(text)
				}
				SendMessage(userID, fmt.Sprintf(MessageDelChannelOK, text))
			}
			SetState(userID, StateIdle)

		} else {
			SendMessage(userID, MessageUnknownCommand)
		}
	}
}

func ProcessCallback(userID int64, messageID int, data string, text string) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
			SendToAdmins(fmt.Sprintf("CallbackError: %s", err))
		}
	}()

	var label int8

	if data[0] == '1' {
		label = 1
	} else {
		label = 0
	}

	text = text[:strings.LastIndex(text, "\n")]

	DisableInlineKeyboard(userID, messageID)
	go ApiTrain(userID, data[1:], text, label)
}

package main

import (
	"log"

	fasthttp "github.com/valyala/fasthttp"
)

type RequestHandler func([]byte)

var BotWebhookPath string
var BotHandler RequestHandler

var PostWebhookPath string
var PostHandler RequestHandler

func InitServer(config ServerConfig, bot_handler RequestHandler, post_handler RequestHandler) {
	BotWebhookPath = "/" + config.Token + "/"
	PostWebhookPath = "/" + config.PostPath + "/"

	BotHandler = bot_handler
	PostHandler = post_handler

	log.Println("Server: inited")
}

func ServerHandler(request *fasthttp.RequestCtx) {
	if string(request.Path()) == BotWebhookPath {
		BotHandler(request.PostBody())
	} else if string(request.Path()) == PostWebhookPath {
		PostHandler(request.PostBody())
	}
}

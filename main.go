package main

import (
	"log"

	fasthttp "github.com/valyala/fasthttp"
)

func main() {

	config := LoadConfig("config.json")

	InitDatabase(config)
	InitApi(config)
	InitDialog(config)
	InitRedis(config)
	InitTelegram(config)

	InitBot(config, ProcessDialog)
	InitServer(config, ProcessRequest, ProcessPost)

	err := fasthttp.ListenAndServe(":"+config.Port, ServerHandler)
	log.Fatal(err)
}

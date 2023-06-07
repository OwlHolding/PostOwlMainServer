import api
import parser
import server
import bot

import json
import logging
import threading


logging.basicConfig(level=logging.INFO, format="%(levelname)s:     %(asctime)s - %(message)s", filename="log.txt",
                    filemode="w")

with open("config.json", 'rb') as file:
    config = json.load(file)

logging.info("Config loaded")

Parser = parser.Parser(config)
API = api.API(config)
Bot = bot.Bot(config)

logging.info("Setting handlers")


async def handler(channel, text):
    users = await API.predict(text, channel, [2009584602])
    Bot.send_post([2009584602], text, channel)


def bot_handler(data):
    Bot.handler(data)


Server = server.Server(config, bot_handler)

logging.info("Finished inited process")

logging.info("Starting threading")

server_thread = threading.Thread(target=Server.start)
server_thread.start()

Parser.start(handler)

logging.info("Finishing process")

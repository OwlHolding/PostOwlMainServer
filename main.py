import api
import parser
import server
import database
import manager
import states

import json
import logging
import threading


logging.basicConfig(level=logging.INFO, format="%(levelname)s:     %(asctime)s - %(message)s", filename="log.txt",
                    filemode="w")

with open("config.json", 'rb') as file:
    config = json.load(file)


bot_states = ['adding-channel', 'deleting-channel']

logging.info("Config loaded")

Parser = parser.Parser(config)
API = api.API(config)
DataBase = database.DataBase(config)
StateStorage = states.StateStorage(config, state_alias=bot_states)

controller = manager.Manager(config, api=API, server=None, database=DataBase, parser=Parser, states=StateStorage)

Server = server.Server(config, controller.bot_handler)

controller.server = Server

logging.info("Finished inited process")

logging.info("Starting threading")

server_thread = threading.Thread(target=Server.start)
server_thread.start()

Parser.start(controller.parser_handler)

logging.info("Finishing process")

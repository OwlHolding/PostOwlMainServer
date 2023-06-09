import logging
import telebot


class Manager:
    """Класс для управления взаимодействием с пользователями"""

    def __init__(self, config, api, server, database, parser, states):
        self.api = api
        self.server = server
        self.database = database
        self.parser = parser
        self.states = states

        self.bot = telebot.TeleBot(config['bot_token'])

        @self.bot.message_handler(content_types=['text'])
        def get_text_messages(message):
            self.bot.send_message(message.from_user.id, "Привет, чем я могу тебе помочь?")

        logging.info("Manager: inited")

    def bot_handler(self, data):
        update = telebot.types.Update.de_json(data)
        self.bot.process_new_updates([update])

    def send_post(self, users: list[int], text: str, channel: str):
        post = f'{text} \n\n @{channel}'
        for user_id in users:
            self.bot.send_message(user_id, post)

    async def parser_handler(self, channel, text):
        print(channel, text)



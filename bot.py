import telebot
import logging


class Bot:

    def __init__(self, config):
        self.bot = telebot.TeleBot(config['bot_token'])

        @self.bot.message_handler(content_types=['text'])
        def get_text_messages(message):
            self.bot.send_message(message.from_user.id, "Привет, чем я могу тебе помочь?")

        logging.info("Bot: inited")

    def handler(self, data):
        update = telebot.types.Update.de_json(data)
        self.bot.process_new_updates([update])

    def send_post(self, users: list[int], text: str, channel: str):
        post = f'{text} \n\n @{channel}'

        for user_id in users:
            self.bot.send_message(user_id, post)

import logging
import telebot
import dialog
import utils
import asyncio


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
            if message.text == '/start':
                self.bot.send_message(message.from_user.id, dialog.HELLO_MESSAGE, parse_mode="HTML")
                self.api.add_user(message.from_user.id)
                self.database.add_user(message.from_user.id)

            elif message.text == '/addchannel':
                self.states.set_state(message.from_user.id, "adding-channel")
                self.bot.send_message(message.from_user.id, dialog.ADDING_CHANNEL, parse_mode="HTML")

            elif message.text == '/delchannel':
                self.states.set_state(message.from_user.id, "deleting-channel")
                self.bot.send_message(message.from_user.id, dialog.DELETING_CHANNEL, parse_mode="HTML")

            elif message.text == '/info':
                data = self.database.get_info(message.from_user.id)
                text = dialog.INFO
                for channel in data:
                    text += f"<code>{channel}</code>\n"
                self.bot.send_message(message.from_user.id, text, parse_mode="HTML")
                self.states.set_state(message.from_user.id, "idle")

            elif message.text == '/cancel':
                self.states.set_state(message.from_user.id, "idle")
                self.bot.send_message(message.from_user.id, dialog.CANCELED, parse_mode="HTML")

            else:
                state = self.states.get_state(message.from_user.id)

                if state == 'adding-channel':
                    if "/" in message.text:
                        channel_name = message.text[message.text.rfind("/")+1:]
                    else:
                        channel_name = message.text

                    if utils.channel_is_exist(channel_name):
                        self.api.add_channel(message.from_user.id, channel_name)
                        self.database.add_channel(message.from_user.id, channel_name)
                        self.parser.add_channel(channel_name)
                        self.bot.send_message(message.from_user.id, dialog.CHANNEL_ADDED, parse_mode="HTML")
                    else:
                        self.bot.send_message(message.from_user.id, dialog.ADDING_UNKNOWN_CHANNEL, parse_mode="HTML")

                    self.states.set_state(message.from_user.id, 'idle')

                elif state == 'deleting-channel':
                    if "/" in message.text:
                        channel_name = message.text[message.text.rfind("/") + 1:]
                    else:
                        channel_name = message.text

                    if channel_name in self.database.get_info(message.from_user.id):
                        self.api.del_channel(message.from_user.id, channel_name)
                        self.database.del_channel(message.from_user.id, channel_name)
                        self.bot.send_message(message.from_user.id, dialog.CHANNEL_DELETED, parse_mode="HTML")
                    else:
                        self.bot.send_message(message.from_user.id, dialog.DELETING_UNKNOWN_CHANNEL, parse_mode="HTML")

                    self.states.set_state(message.from_user.id, 'idle')

                else:
                    self.bot.send_message(message.from_user.id, dialog.UNKNOWN_COMMAND, parse_mode="HTML")

        self.bot.set_webhook(config['webhook_url'])
        logging.info("Manager: inited")

    def bot_handler(self, data):
        update = telebot.types.Update.de_json(data)
        self.bot.process_new_updates([update])

    def send_post(self, users: list[int], text: str, channel: str):
        post = f'{text} \n\n @{channel}'
        for user_id in users:
            self.bot.send_message(user_id, post)

    async def parser_handler(self, channel, text):
        users = self.database.get_users(channel)
        users = self.api.predict(text, channel, users)
        self.send_post(users, text, channel)


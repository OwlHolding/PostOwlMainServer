from telethon import TelegramClient, events
import logging


class Parser:
    """Класс для работы с telegram-каналами"""

    def __init__(self, config):
        self.client = None
        self.config = config
        logging.info("TelegramParser inited")

    def start(self, handler):

        self.client = TelegramClient(session="session",
                                     api_id=self.config['telegram_api_id'],
                                     api_hash=self.config['telegram_api_hash'],
                                     system_version="4.16.30-vxCUSTOM")

        with self.client:

            @self.client.on(events.NewMessage())
            async def func(event):
                message = event.message.message
                channel = await self.client.get_entity(event.message.peer_id)
                channel_name = channel.username
                if channel_name and not ("bot" in channel_name.lower()):
                    logging.info(f"Parser: Getting post from channel {channel_name}")
                    await handler(channel_name, message)

            logging.info("Parser: starting listening")
            self.client.run_until_disconnected()

from telethon import TelegramClient, events
from telethon.tl.functions.channels import JoinChannelRequest
import logging
import asyncio
import queue


class Parser:
    """Класс для работы с telegram-каналами"""

    channel_queue = queue.Queue()

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

            async def checker():
                while True:
                    if self.channel_queue.qsize() > 0:
                        channel_name = self.channel_queue.get()

                        channel_entity = await self.client.get_entity(channel_name)
                        await self.client(JoinChannelRequest(channel_entity))
                        logging.info(f"Parser: joined to channel {channel_name}")
                    await asyncio.sleep(1)

            logging.info("Parser: starting listening")

            loop = asyncio.get_event_loop()
            loop.create_task(checker())

            self.client.run_until_disconnected()
            loop.close()

    def add_channel(self, channel: str):
        self.channel_queue.put(channel)

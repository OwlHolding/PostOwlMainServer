import json
import requests
import logging


class API:
    """Класс для общения с ML-серверами. Реализовывает распределение нагрузки на систему."""

    def __init__(self, config: dict):
        """Принимает список IP адресов ML-серверов"""
        if len(config['ml_api_ips']) == 1:
            self.ips = config['ml_api_ips'][0]
        else:
            self.ips = {}
            for key in config['ml_api_ips']:
                self.ips[key] = 0

        logging.info("API: inited")

    def get_ip(self, value: int):
        if len(self.ips) == 1:
            return self.ips

        else:
            ip = sorted(self.ips, key=lambda x: self.ips[x])[0]
            self.ips[ip] = value
            return ip

    async def add_user(self, user_id: int):
        logging.info(f"API: add user {user_id}")
        requests.post(f"http://{self.get_ip(1)}/add-user/{user_id}/")

    async def del_user(self, user_id: int):
        logging.info(f"API: del user {user_id}")
        requests.delete(f"http://{self.get_ip(1)}/del-user/{user_id}/")

    async def add_channel(self, user_id: int, channel: str):
        logging.info(f"API: add channel {channel} for {user_id}")
        requests.post(f"http://{self.get_ip(1)}/add-channel/{user_id}/{channel}/")

    async def del_channel(self, user_id: int, channel: str):
        logging.info(f"API: del channel {channel} for {user_id}")
        requests.delete(f"http://{self.get_ip(1)}/del-channel/{user_id}/{channel}/")

    async def predict(self, post: str, channel: str, users: list[int]) -> list:
        logging.info(f"API: predict request for {channel} and users: {users}")
        data = json.dumps({
            "post": post,
            "channel": channel,
            "users": users
        })

        return json.loads(requests.post(f"http://{self.get_ip(4)}/predict/", data=data).content)

    async def train(self, user_id: int, channel: str, post: str, label: bool):
        logging.info(f"API: train request for {channel} from {user_id}")
        data = json.dumps({
            "text": post,
            "label": label
        })

        requests.put(f"http://{self.get_ip(5)}/train/{user_id}/{channel}/", data=data)

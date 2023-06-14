import json
import logging
from requests_futures.sessions import FuturesSession


class API:
    """Класс для общения с ML-серверами. Реализовывает распределение нагрузки на систему."""

    def __init__(self, config: dict):
        """Принимает список IP адресов ML-серверов"""

        if len(config['ml_api_ips']) == 1:
            self.ips = config['ml_api_ips']
        else:
            self.ips = {}
            for key in config['ml_api_ips']:
                self.ips[key] = 0

        self.session = FuturesSession()

        logging.info("API: inited")

    def get_ip(self, value: int):
        if len(self.ips) == 1:
            return self.ips[0]

        else:
            print(self.ips)
            ip = sorted(self.ips, key=lambda x: self.ips[x])[0]
            self.ips[ip] = value
            return ip

    def add_user(self, user_id: int):
        logging.info(f"API: add user {user_id}")
        self.session.post(f"http://{self.get_ip(1)}/add-user/{user_id}/")

    def del_user(self, user_id: int):
        logging.info(f"API: del user {user_id}")
        self.session.delete(f"http://{self.get_ip(1)}/del-user/{user_id}/")

    def add_channel(self, user_id: int, channel: str):
        logging.info(f"API: add channel {channel} for {user_id}")
        self.session.post(f"http://{self.get_ip(1)}/add-channel/{user_id}/{channel}/")

    def del_channel(self, user_id: int, channel: str):
        logging.info(f"API: del channel {channel} for {user_id}")
        self.session.delete(f"http://{self.get_ip(1)}/del-channel/{user_id}/{channel}/")

    def predict(self, post: str, channel: str, users: list[int]) -> list:
        logging.info(f"API: predict request for {channel} and users: {users}")
        data = json.dumps({
            "post": post,
            "channel": channel,
            "users": users
        })

        request = self.session.post(f"http://{self.get_ip(4)}/predict/", data=data)
        response = request.result()

        return json.loads(response.content)['users']

    def train(self, user_id: int, channel: str, post: str, label: bool):
        logging.info(f"API: train request for {channel} from {user_id}")
        data = json.dumps({
            "text": post,
            "label": label
        })

        self.session.put(f"http://{self.get_ip(5)}/train/{user_id}/{channel}/", data=data)

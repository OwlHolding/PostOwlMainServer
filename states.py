import json
import redis
import logging


class StateStorage:
    """Класс для хранения состояний конечного автомата"""

    def __init__(self, config, state_alias):
        self.redis = redis.Redis(
            host=config['redis_host'],
            port=config['redis_port'],
            decode_responses=True
        )
        try:
            self.redis.ping()
        except redis.exceptions.ConnectionError:
            logging.info("StateStorage: init failed, can not connect to Redis")
        else:
            logging.info("StateStorage: inited")

        self.state_alias = state_alias

    def get_state(self, user_id: int):
        try:
            index = self.redis.get(user_id)
            return self.state_alias[int(index)]
        except redis.exceptions.ConnectionError:
            logging.info('StateStorage: can not connect to Redis')

    def set_state(self, user_id: int, state) -> bool:
        try:
            index = self.state_alias.index(state)
            return self.redis.set(user_id,  index, 24 * 60 * 60)
        except redis.exceptions.ConnectionError:
            logging.info('StateStorage: can not connect to Redis')

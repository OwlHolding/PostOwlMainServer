import json
import redis
import logging


class StateStorage:

    def __init__(self, config):
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

    def get_state(self, user_id: int):
        try:
            return self.redis.get(user_id)
        except redis.exceptions.ConnectionError:
            logging.info('StateStorage: can not connect to Redis')

    def set_state(self, user_id: int, state) -> bool:
        try:
            return self.redis.set(user_id, state)
        except redis.exceptions.ConnectionError:
            logging.info('StateStorage: can not connect to Redis')

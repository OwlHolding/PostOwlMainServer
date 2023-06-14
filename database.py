import mysql.connector
from mysql.connector import Error
import logging
import json


class DataBase:
    """Класс для работы с базой данных"""

    def __init__(self, config):

        self.connection = mysql.connector.connect(
            user=config['db_user'],
            password=config['db_password'],
            host=config['db_host'])

        self.cursor = self.connection.cursor()

        logging.info("DataBase: creating database")
        try:
            self.cursor.execute("CREATE DATABASE postowl")
            self.connection.commit()
        except Error:
            logging.info("DataBase: database already exist")

        self.connection = mysql.connector.connect(
            user=config['db_user'],
            password=config['db_password'],
            host=config['db_host'],
            database="postowl"
        )
        self.cursor = self.connection.cursor()

        query = "CREATE TABLE users(id BIGINT PRIMARY KEY, channels TEXT NOT NULL DEFAULT '')"
        logging.info("DataBase: creating table")
        try:
            self.cursor.execute(query)
            self.connection.commit()
        except Error:
            logging.info("DataBase: table already exist")

        logging.info("DataBase: inited")

    def execute_query(self, query: str) -> bool:
        try:
            self.cursor.execute(query)
            self.connection.commit()
        except Error as e:
            logging.info(f"DataBase: Failed to execute query: {e}")
            return False
        return True

    def execute_read_query(self, query: str):
        try:
            self.cursor.execute(query)
            result = self.cursor.fetchall()

            return result
        except Error as e:
            logging.info(f"DataBase: Failed to execute query: {e}")

    def add_user(self, user_id: int):
        query = f"INSERT INTO users VALUES ({user_id}, '')"
        self.execute_query(query)

    def del_user(self, user_id: int):
        query = f"DELETE FROM users WHERE id = {user_id}"
        self.execute_query(query)

    def add_channel(self, user_id: int, channel: str):
        query = f"SELECT * FROM users WHERE id = {user_id}"
        channels = self.execute_read_query(query)[0][1].split("&")
        if channels[0] == '':
            channels[0] = channel
        else:
            if not channel in channels:
                channels.append(channel)

        query = f"UPDATE users SET channels = '{'&'.join(channels)}' WHERE id = {user_id}"
        self.execute_query(query)

    def del_channel(self, user_id: int, channel: str):
        query = f"SELECT * FROM users WHERE id = {user_id}"
        channels = self.execute_read_query(query)[0][1].split("&")

        try:
            channels.pop(channels.index(channel))
        except ValueError:
            pass
        else:
            query = f"UPDATE users SET channels = '{'&'.join(channels)}' WHERE id = {user_id}"
            self.execute_query(query)

    def get_users(self, channel: str) -> list[int]:
        query = f"SELECT id FROM users WHERE channels LIKE '%{channel}%'"
        users = []
        for user in self.execute_read_query(query):
            users.append(user[0])
        return users

    def get_info(self, user_id: int):
        query = f"SELECT channels FROM users WHERE id = {user_id}"
        channels = self.execute_read_query(query)[0][0].split('&')
        return channels


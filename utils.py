import requests


def channel_is_exist(channel):
    """Проверка на канала"""
    resp = requests.get(f"https://rsshub.app/telegram/channel/{channel}")
    if resp.status_code != 200:
        return False
    return True

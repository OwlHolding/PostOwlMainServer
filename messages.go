package main

const (
	MessageHello                = "Привет!\nЯ – Почтовая сова. Моя главная задача – экономить твоё время. Я обучаюсь на твоих предпочтениях и каждый день присылаю сводку самых полезных постов за сутки. Для начала нам нужно добавить хотя бы один канал с помощью команды /addchannel, и установить время рассылки командой /changetime"
	MessageAddChannel           = "Пришли мне название канала или ссылку на него. Например, <code>forbesrussia</code>"
	MessageChannelAlreadyAdded  = "Такой канал уже добавлен ранее"
	MessageChannelNotExists     = "К сожалению, я не могу найти такой канал 😭"
	MessageAddChannelOK         = "Канал <code>%s</code> успешно добавлен"
	MessageChannelOverflow      = "К сожалению, пока нельзя добавить более %d каналов 😭"
	MessageRateCycleEnd         = "Обучение завершено, спасибо 🙂"
	MessageRateCycleFormat      = "Пришли 👍 или 👎, пожалуйста"
	MessageRateCycleWait        = "Нужно немного времени, чтобы модель обучилась"
	MessageRateCycleAllPositive = "Нужен хотя бы один отрицательный пример. Перешли, пожалуйста, один пост, который тебе неинтересен"
	MessageRateCycleAllNegative = "Нужен хотя бы один положительный пример. Перешли, пожалуйста, один пост, который тебе интересен"
	MessageDelChannel           = "Пришли мне название канала или ссылку на канал, который нужно удалить"
	MessageChannelNotListed     = "Такого канал нет в списке рассылки"
	MessageDelChannelOK         = "Канал <code>%s</code> успешно удален"
	MessageChangeTime           = "Пришли мне время рассылки, которое тебе удобно. Например, 17:00"
	MessageTimeInvalidFormat    = "Пришли время в формате часы:минуты, пожалуйста"
	MessageChangeTimeOK         = "Спасибо! Теперь я буду прислать рассылку в %s"
	MessageUserDisabled         = "Рассылка успешно отключена. Для повторной активации заново установи время с помощью команды /changetime"
	MessageInfo                 = "Добавленные каналы:\n\n%s"
	MessageNoNewPosts           = "Интересных постов сегодня нет 😒"
	MessageDailySummary         = "Ежедневная сводка постов на сегодня☺️"
	MessageNotForget            = "Не забудь назначить время рассылки командой /changetime 🙂"
	MessageUnknownCommand       = "Неизвестная команда"
	MessageCancel               = "Команда отменена"
	MessageError                = "Похоже, что вам удалось найти ошибку в нашем коде! Большое спасибо, мы уже заняты ее устранением"
	MessageBanned               = "Вы плохо себя вели и были забанены 😔"
	MessageNotAllowed           = "На данный момент проект PostOwl функционирует в демонстрационном режиме и доступен ограниченному кругу лиц. Пожалуйста, пришлите ключ, если он был предоставлен Вам разработчиками"
)
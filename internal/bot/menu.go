package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func createMenu() tgbotapi.ReplyKeyboardMarkup {
	var TimerButtonText, AddSetButtonText, SkipSetButton string
	if isRunning {
		TimerButtonText = "Завершить тренировку"
	} else {
		TimerButtonText = "Начать тренировку"
	}

	AddSetButtonText = "Занести подход"
	SkipSetButton = "Пропустить подход"

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(AddSetButtonText), tgbotapi.NewKeyboardButton(SkipSetButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(TimerButtonText),
		),
	)

	return keyboard
}

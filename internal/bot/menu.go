package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	AddSetButtonText                   string = "Добавить подход"
	NextExerciseButtonText             string = "Следующее упражнение"
	NewSessionButtonText               string = "Новая тренировка"
	PreviousSessionsButtonText         string = "Прошлые тренировки"
	StatisticsButtonText               string = "Статистика"
	IndividualDataButtonText           string = "Индивидуальные данные"
	EditSessionButtonText              string = "Редактировать тренировку"
	EditPreviousSessionButtonText      string = "Редактировать завершенную тренировку"
	MainMenuButtonText                 string = "Главное меню"
	AllTimeStatisticsButtonText        string = "Статистика за все время"
	OverTimePeriodStatisticsButtonText string = "Статистика за период"
	ShowLastSessionsButtonText         string = "Показать последние 10 тренировок"
)

var (
	TimerButtonText string = "Начать тренировку"
	CurrentPage     string = "Main menu" // Main menu - Current session - Previous sessions - Statistics
)

func createMenu() tgbotapi.ReplyKeyboardMarkup {

	if timerRunning {
		TimerButtonText = "Завершить тренировку"
	} else {
		TimerButtonText = "Начать тренировку"
	}

	var keyboard tgbotapi.ReplyKeyboardMarkup

	switch CurrentPage {
	case "Main menu":
		{
			keyboard = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(NewSessionButtonText), tgbotapi.NewKeyboardButton(PreviousSessionsButtonText)),
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(StatisticsButtonText)))
		}
	case "Current session":
		{
			keyboard = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(AddSetButtonText), tgbotapi.NewKeyboardButton(NextExerciseButtonText)),
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(IndividualDataButtonText), tgbotapi.NewKeyboardButton(EditSessionButtonText)),
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(MainMenuButtonText), tgbotapi.NewKeyboardButton(TimerButtonText)),
			)
		}
	case "Previous sessions":
		{
			keyboard = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(ShowLastSessionsButtonText), tgbotapi.NewKeyboardButton(EditPreviousSessionButtonText)),
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(MainMenuButtonText)),
			)
		}
	case "Statistics":
		{
			keyboard = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(AllTimeStatisticsButtonText), tgbotapi.NewKeyboardButton(OverTimePeriodStatisticsButtonText)),
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(MainMenuButtonText)),
			)
		}
	default:
		{
			keyboard = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(NewSessionButtonText), tgbotapi.NewKeyboardButton(PreviousSessionsButtonText)),
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(StatisticsButtonText)))
		}
	}

	return keyboard
}

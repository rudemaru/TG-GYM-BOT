package bot

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	startTime time.Time
	isRunning bool
)

type Bot struct {
	API *tgbotapi.BotAPI
}

func NewBot(token string) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("Bot initialization error: %w", err)
	}

	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	return &Bot{API: botAPI}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		b.handleCommand(update.Message)
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Text {
	case "Начать тренировку":
		b.startTimer(message)
	case "Завершить тренировку":
		b.stopTimer(message)
	default:
		b.sendDefaultMessage(message)
	}
}

func (b *Bot) startTimer(message *tgbotapi.Message) {
	if !isRunning {
		startTime = time.Now()
		isRunning = true
		msg := tgbotapi.NewMessage(message.Chat.ID, "Таймер запущен.")
		msg.ReplyMarkup = createMenu()
		b.API.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Таймер уже запущен.")
		msg.ReplyMarkup = createMenu()
		b.API.Send(msg)
	}
}

func (b *Bot) stopTimer(message *tgbotapi.Message) {
	if isRunning {
		duration := time.Since(startTime)
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Таймер остановлен. Длительность тренировки: %s", duration))
		isRunning = false
		msg.ReplyMarkup = createMenu()
		b.API.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Таймер не запущен.")
		msg.ReplyMarkup = createMenu()
		b.API.Send(msg)
	}
}

func (b *Bot) sendDefaultMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Исполльзуйте меню ниже для управления ботом.")
	msg.ReplyMarkup = createMenu()
	b.API.Send(msg)
}

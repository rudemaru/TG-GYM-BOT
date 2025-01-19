package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	startTime      time.Time
	timerRunning   bool
	weightChoosing bool
	repChoosing    bool
	workingSet     *set
)

type Bot struct {
	API *tgbotapi.BotAPI
}

type set struct {
	weight float32
	reps   int
	valid  bool
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
		{
			b.startTimer(message)
		}
	case "Завершить тренировку":
		{
			b.stopTimer(message)
		}
	case "Добавить подход":
		{
			//log.Printf("[handleCommand]WeightChoosing: %v, RepChoosing: %v, WorkingSet is nil: %v", weightChoosing, repChoosing, workingSet == nil)
			b.addSet(message)
		}
	case "Новая тренировка":
		{
			CurrentPage = "Current session"
			b.sendResponse(message, "Новая тренировка")
		}
	case "Прошлые тренировки":
		{
			CurrentPage = "Previous sessions"
			b.sendResponse(message, "Прошлые тренировки")
		}
	case "Статистика":
		{
			CurrentPage = "Statistics"
			b.sendResponse(message, "Статистика")
		}
	case "Главное меню":
		{
			CurrentPage = "Main menu"
			b.sendResponse(message, "Главное меню")
		}
	default:
		b.handleDefault(message)
	}
}

func (b *Bot) handleDefault(message *tgbotapi.Message) {
	if weightChoosing {
		b.processWeightInput(message, workingSet)
	} else if repChoosing {
		b.processRepetitionsInput(message, workingSet)
	} else {
		b.sendDefaultMessage(message)
	}
}

func (b *Bot) startTimer(message *tgbotapi.Message) {
	if !timerRunning {
		startTime = time.Now()
		timerRunning = true
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
	if !timerRunning {
		b.sendResponse(message, "Таймер не запущен.")
		return
	}

	duration := time.Since(startTime)
	durationFormatted := formatDuration(duration)

	msg := fmt.Sprintf("Таймер остановлен. Длительность тренировки: %s", durationFormatted)
	timerRunning = false
	b.sendResponse(message, msg)

	if workingSet != nil {
		workingSet.valid = false
	}
}

func (b *Bot) addSet(message *tgbotapi.Message) {
	if !timerRunning {
		b.sendResponse(message, "Тренировка не начата.")
		return
	}

	//log.Printf("[addSet]WeightChoosing: %v, RepChoosing: %v, WorkingSet is nil: %v", weightChoosing, repChoosing, workingSet == nil)

	if workingSet == nil || !workingSet.valid {
		workingSet = new(set)
		workingSet.valid = true
		b.promptWeightInput(message, workingSet)
	} else if weightChoosing {
		b.processWeightInput(message, workingSet)
	} else if repChoosing {
		b.processRepetitionsInput(message, workingSet)
	} else {
		b.promptRepetitionsInput(message, workingSet)
	}
}

func (b *Bot) sendDefaultMessage(message *tgbotapi.Message) {
	b.sendResponse(message, "Это еще не реализовано :(")
}

func (b *Bot) sendResponse(message *tgbotapi.Message, text string) {
	//log.Printf("[sendResponse]WeightChoosing: %v, RepChoosing: %v, WorkingSet is nil: %v", weightChoosing, repChoosing, workingSet == nil)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = createMenu()
	b.API.Send(msg)
}

func (b *Bot) promptWeightInput(message *tgbotapi.Message, workingSet *set) {
	msg := "Введите вес в килограммах."
	weightChoosing = true
	repChoosing = false
	b.sendResponse(message, msg)
}

func (b *Bot) processWeightInput(message *tgbotapi.Message, workingSet *set) {
	inputWeight := strings.ReplaceAll(message.Text, ",", ".")
	parsedWeight, err := strconv.ParseFloat(inputWeight, 32)

	if err != nil {
		b.sendResponse(message, "Введите действительное число для веса.")
		return
	}

	workingSet.weight = float32(parsedWeight)
	weightChoosing = false
	b.promptRepetitionsInput(message, workingSet)
}

func (b *Bot) promptRepetitionsInput(message *tgbotapi.Message, workingSet *set) {
	msg := "Введите количество повторов."
	repChoosing = true
	b.sendResponse(message, msg)
}

func (b *Bot) processRepetitionsInput(message *tgbotapi.Message, workingSet *set) {
	inputReps, err := strconv.ParseInt(message.Text, 0, 32)

	if err != nil {
		b.sendResponse(message, "Введите действительное число для повторов.")
		return
	}

	workingSet.reps = int(inputReps)
	msg := fmt.Sprintf("Подход %vкг * %v успешно сохранен.", workingSet.weight, workingSet.reps)
	b.sendResponse(message, msg)

	workingSet.valid = false
	weightChoosing = false
	repChoosing = false
}

func formatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02dч %02dм %02dс", hours, minutes, seconds)
}

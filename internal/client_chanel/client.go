package clientchanel

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"strings"
	"sync"
	"tech-sup-bot-pl/internal/keyboards"
	"tech-sup-bot-pl/tgUtil"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	// "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type UpdateHandlerFunc func(update tgbotapi.Update, bot *tgbotapi.BotAPI)

type Beluga struct {
	Logger *zap.Logger
	State  map[int64]UpdateHandlerFunc
	Sword  sync.Mutex //rw mutex
}

func New() *Beluga {
	bl := &Beluga{
		State: make(map[int64]UpdateHandlerFunc),
		Sword: sync.Mutex{},
	}
	return bl
}

func (b *Beluga) Get(chadID int64) (f UpdateHandlerFunc, ok bool) {
	// b.Sword.Lock()
	f, ok = b.State[chadID]
	// b.Sword.Unlock()
	return f, ok
}

func (b *Beluga) Add(chadID int64, f func(update tgbotapi.Update, bot *tgbotapi.BotAPI)) {
	// b.Sword.Lock()
	b.State[chadID] = f
	// b.Sword.Unlock()
}

func createSignature(inputData string) string {
	if strings.Contains(inputData, ",") {
		str := strings.Split(inputData, ",")
		key := str[0]
		message := ""
		for i := 1; i < len(str); i++ {
			if i == (len(str) - 1) {
				message += str[i]
			} else {
				message += str[i] + ","
			}
		}
		// logrus.Printf(key)
		// logrus.Printf(message)
		signature := hmac.New(sha512.New, []byte(key))
		signature.Write([]byte(message))
		return hex.EncodeToString(signature.Sum(nil))
	} else {
		return "Ошибка в заполнении входных данных"
	}
}

func (b *Beluga) DigitalSignature(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.CallbackQuery != nil {

		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := bot.Request(callback); err != nil {
			// logrus.Error(err)
		}
		switch update.CallbackQuery.Data {
		case "Меню":
			// logrus.Printf(update.CallbackQuery.Data)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.MenuReply)
			msg.ReplyMarkup = keyboards.StartKeyBoard
			tgUtil.SendBotMessage(msg, bot)
			delete(b.State, update.CallbackQuery.Message.Chat.ID)
		}

	} else {

		txt := update.Message.Text
		// log.Printf(txt)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваша подпись: \n \n"+createSignature(txt))
		msg.ReplyMarkup = keyboards.ToMainTheme
		tgUtil.SendBotMessage(msg, bot)

	}
}

func (b *Beluga) Feedback(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.CallbackQuery != nil {

		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := bot.Request(callback); err != nil {
			// logrus.Error(err)
		}
		switch update.CallbackQuery.Data {
		case "Меню":
			// logrus.Printf(update.CallbackQuery.Data)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.MenuReply)
			msg.ReplyMarkup = keyboards.StartKeyBoard
			tgUtil.SendBotMessage(msg, bot)
			delete(b.State, update.CallbackQuery.Message.Chat.ID)
		}

	} else {
		txt := "ХР\n" + update.Message.Chat.FirstName + " " + update.Message.Chat.LastName + "\n" + "@" + update.Message.Chat.UserName + "\n" + update.Message.Text
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyboards.GratitudeReply)
		msg.ReplyMarkup = keyboards.ToMainTheme
		tgUtil.SendBotMessage(msg, bot)
		msg = tgbotapi.NewMessageToChannel("1661385575", txt) // канал вынести в отдельную переменную окружения
		// tgUtil.SendBotMessage(msg, bot)
	}
}

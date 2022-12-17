package clientchanel

import (
	"log"
	"sync"
	"tech-sup-bot-pl/internal/keyboards"
	"tech-sup-bot-pl/tgUtil"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Beluga struct {
	state map[int64]struct{}
	sword sync.Mutex
}

func New() *Beluga {
	bl := &Beluga{
		state: make(map[int64]struct{}),
		sword: sync.Mutex{},
	}

	return bl
}

func (b *Beluga) Check(chadID int64) bool {
	b.sword.Lock()
	defer b.sword.Unlock()

	if _, ok := b.state[chadID]; ok {
		return true
	}

	return false
}

func (b *Beluga) Cycle(updatesChan tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) {
	for newUpdate := range updatesChan {
		if newUpdate.CallbackQuery != nil {

			callback := tgbotapi.NewCallback(newUpdate.CallbackQuery.ID, newUpdate.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				logrus.Error(err)
			}
			switch newUpdate.CallbackQuery.Data {
			case "Меню":
				logrus.Printf(newUpdate.CallbackQuery.Data)
				replay := "Вы в главном меню, выберете на экране интересующий вас раздел"
				msg := tgbotapi.NewMessage(newUpdate.CallbackQuery.Message.Chat.ID, replay)
				msg.ReplyMarkup = keyboards.StartKeyBoard
				tgUtil.SendBotMessage(msg, bot)
			}
		} else if newUpdate.Message == nil {
			continue
		} else {
			txt := newUpdate.Message.Text
			log.Print(txt)
			msg := tgbotapi.NewMessage(newUpdate.Message.Chat.ID, "Ваша подпись: \n \n")
			msg.ReplyMarkup = keyboards.ToMainTheme
			tgUtil.SendBotMessage(msg, bot)
			break
		}
		break
	}
}

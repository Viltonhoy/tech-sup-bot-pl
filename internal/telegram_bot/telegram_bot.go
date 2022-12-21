package telegram_bot

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	//"github.com/joho/godotenv"

	tg_bot_for_ts "tech-sup-bot-pl/faq"
	clientchanel "tech-sup-bot-pl/internal/client_chanel"
	"tech-sup-bot-pl/internal/keyboards"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	//"tg-bot-for-ts/repository"
	//"tg-bot-for-ts/service"
	"tech-sup-bot-pl/tgUtil"
)

type Bot struct {
	Logger  *zap.Logger
	Bot     *tgbotapi.BotAPI
	Updates tgbotapi.UpdatesChannel
}

func NewBot(logger *zap.Logger) (*Bot, error) {
	if logger == nil {
		return nil, errors.New("no logger provided")
	}

	bot, err := tgbotapi.NewBotAPI("5957121655:AAHQweDAwIECC_ppdAYEhOG3vsg0cAwQaN8")
	if err != nil {
		logger.Error("failed to create a new BotAPI instance", zap.Error(err))
		return nil, err
	}
	bot.Debug = true

	logger.Debug("Authorized on ", zap.String("account", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	return &Bot{
		Bot:     bot,
		Updates: updates,
	}, err
}

func (b *Bot) sendBotMessage(msg tgbotapi.MessageConfig) {
	if _, err := b.Bot.Send(msg); err != nil {
		b.Logger.Error("", zap.Error(err))
	}
}

func (b *Bot) BotWorker(cc *clientchanel.Beluga) error {
	for update := range b.Updates {

		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		var chatID int64

		if update.Message != nil {
			chatID = update.Message.Chat.ID
		}
		if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
		}

		if f, ok := cc.Get(chatID); ok {
			f(update, b.Bot)
			continue
		}

		if update.Message != nil {
			b.switcherMessage(update)
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := b.Bot.Request(callback); err != nil {
				b.Logger.Error("sending Chattable error", zap.Error(err))
				return err
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			b.switcherCollback(update, msg, cc)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyboards.UnrecognizedСommand)
			tgUtil.SendBotMessage(msg, b.Bot)
		}
	}
	return nil
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func (b *Bot) switcherMessage(update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyboards.StartReply)
		msg.ReplyMarkup = keyboards.StartKeyBoard

		tgUtil.SendBotMessage(msg, b.Bot)

	case "commands":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyboards.CommandsReply)
		tgUtil.SendBotMessage(msg, b.Bot)

	case "main":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyboards.MainReply)
		msg.ReplyMarkup = keyboards.StartKeyBoard
		tgUtil.SendBotMessage(msg, b.Bot)

	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyboards.DefReply)
		tgUtil.SendBotMessage(msg, b.Bot)
	}
}

func (b *Bot) switcherCollback(update tgbotapi.Update, msg tgbotapi.MessageConfig, cc *clientchanel.Beluga) {
	switch update.CallbackQuery.Data {

	case "Выберите тип проблемы":
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.TypeReply)
		msg.ReplyMarkup = keyboards.KeyBoard
		tgUtil.SendBotMessage(msg, b.Bot)

	case "Генерация подписи":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.DigitalSignatureGeneration)
		msg.ReplyMarkup = keyboards.ToMainTheme
		tgUtil.SendBotMessage(msg, b.Bot)
		cc.Add(update.CallbackQuery.Message.Chat.ID, cc.DigitalSignature)

	case "Обратная связь":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.FeedbackReply)
		msg.ReplyMarkup = keyboards.ToMainTheme
		tgUtil.SendBotMessage(msg, b.Bot)
		cc.Add(update.CallbackQuery.Message.Chat.ID, cc.Feedback)

	case "FAQ":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.FaqReply)
		msg.ReplyMarkup = keyboards.FAQKeyBoard
		tgUtil.SendBotMessage(msg, b.Bot)

	case "Регистрация":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.Registration)
		msg.ReplyMarkup = keyboards.FAQKeyBoard
		tgUtil.SendBotMessage(msg, b.Bot)

	case "Signature":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.SignatureQuestion)
		msg.ReplyMarkup = keyboards.FAQKeyBoard
		tgUtil.SendBotMessage(msg, b.Bot)

	case "API":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.APIQuestion)
		msg.ReplyMarkup = keyboards.FAQKeyBoard
		tgUtil.SendBotMessage(msg, b.Bot)

	case "IT":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.ITQuestion)
		msg.ReplyMarkup = keyboards.FAQKeyBoard
		tgUtil.SendBotMessage(msg, b.Bot)

	case "Заявки":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.RequestQuestion)
		msg.ReplyMarkup = keyboards.FAQKeyBoard
		tgUtil.SendBotMessage(msg, b.Bot)

	case "Меню":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.MenuReply)
		msg.ReplyMarkup = keyboards.StartKeyBoard
		tgUtil.SendBotMessage(msg, b.Bot)

	default:
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.InvalidReqReply)
		tgUtil.SendBotMessage(msg, b.Bot)
	}
}

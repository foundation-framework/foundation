package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/intale-llc/foundation/notify"
)

type telegram struct {
	bot *tgbotapi.BotAPI
}

func NewProvider(token string) (notify.Provider, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &telegram{
		bot: bot,
	}, nil
}

func (t *telegram) Send(text string, chatIds []int64, attachments ...notify.Attachment) error {
	//defer func() {
	//	for _, attachment := range attachments {
	//		// Ignore any errors at this stage
	//		_ = attachment.Close()
	//	}
	//}()

	for _, chatId := range chatIds {
		if err := t.sendText(chatId, text); err != nil {
			return err
		}

		if err := t.sendAttachments(chatId, attachments); err != nil {
			return err
		}
	}

	return nil
}

func (t *telegram) sendText(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "HTML"

	_, err := t.bot.Send(msg)
	return err
}

func (t *telegram) sendAttachments(chatId int64, attachments []notify.Attachment) error {
	for _, attachment := range attachments {
		if err := attachment.Reset(); err != nil {
			return err
		}

		file := tgbotapi.FileReader{
			Name:   attachment.Name(),
			Reader: attachment.Reader(),
		}

		_, err := t.bot.SendMediaGroup(
			tgbotapi.NewMediaGroup(
				chatId,
				[]interface{}{
					tgbotapi.NewInputMediaDocument(file),
				},
			),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

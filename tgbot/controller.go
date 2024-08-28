package tgbot

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Controller struct {
	Session *Session
	Param   string
	Start   string
	update  tgbotapi.Update
	bot     *tgbotapi.BotAPI
}

type ControllerInterface interface {
	Init(update tgbotapi.Update, bot *tgbotapi.BotAPI, text string)
	Handle()
	HandleEdit()
	HandleNext() bool
}

func (c *Controller) Init(update tgbotapi.Update, b *tgbotapi.BotAPI, text string) {
	c.bot = b
	c.Param = text
	c.update = update
	c.Session = SessionMgr.GetSession(c.GetUserInfo().ID)
}

func (c *Controller) Handle() {
	msg := "bot handler not implemented"
	log.Println(msg)
	c.SendError(errors.New(msg))
}

func (c *Controller) HandleEdit() {
	msg := "bot edit handler not implemented"
	log.Println(msg)
	c.SendError(errors.New(msg))
}

func (c *Controller) HandleNext() bool {
	msg := "bot anyhandler not implemented"
	log.Println(msg)
	c.SendError(errors.New(msg))
	return false
}

func (c *Controller) DefaultMenu() {
	c.setMenuButton("default", "", "")
}

func (c *Controller) ShowMenu(text, url string) {
	c.setMenuButton("web_app", text, url)
}

func (c *Controller) setMenuButton(buttonType, text, url string) {
	_, err := c.bot.Request(tgbotapi.SetChatMenuButtonConfig{
		ChatID: c.ChatId(),
		MenuButton: &tgbotapi.MenuButton{
			Type: buttonType,
			Text: text,
			WebApp: &tgbotapi.WebAppInfo{
				URL: url,
			},
		},
	})
	if err != nil {
		log.Println("SetMenuButton error:", err)
	}
}

func (c *Controller) Reply(text string) {
	msg := tgbotapi.NewMessage(c.ChatId(), text)
	msg.ReplyToMessageID = c.update.Message.MessageID
	c.send(msg)
}

func (c *Controller) SendMsg(text string) {
	msg := tgbotapi.NewMessage(c.ChatId(), text)
	c.send(msg)
}

func (c *Controller) DeleteLastBotMsg() {
	c.deleteMessage(c.Session.LastBotId)
}

func (c *Controller) DeleteLastUserMsg() {
	c.deleteMessage(c.Session.LastUserId)
}

func (c *Controller) deleteMessage(messageId int) {
	if messageId == 0 {
		return
	}
	msg := tgbotapi.NewDeleteMessage(c.ChatId(), messageId)
	c.sendEdit(msg)
}

func (c *Controller) EditLastBotMsg(text string) {
	msg := tgbotapi.NewEditMessageText(c.ChatId(), c.Session.LastBotId, text)
	c.sendEdit(msg)
}

func (c *Controller) EditLastBotMsgWithUrl(text string, buttons [][]Button) {
	msg := tgbotapi.NewEditMessageText(c.ChatId(), c.Session.LastBotId, text)
	msg.ReplyMarkup = c.makeInlineKeyboard(buttons)
	c.sendEdit(msg)
}

func (c *Controller) EditLastBotPhoto(path, caption string) {
	c.editPhotoWithButtons(path, caption, nil)
}

func (c *Controller) EditLastBotPhotoWithUrl(path, caption string, buttons [][]Button) {
	c.editPhotoWithButtons(path, caption, buttons)
}

func (c *Controller) editPhotoWithButtons(path, caption string, buttons [][]Button) {
	buff, err := os.ReadFile(path)
	if err != nil {
		log.Println("file read error:", err)
		return
	}

	file := tgbotapi.FileBytes{Name: "photo", Bytes: buff}
	baseMedia := tgbotapi.BaseInputMedia{
		Type:      "photo",
		Media:     file,
		Caption:   caption,
		ParseMode: "HTML",
	}

	msg := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      c.ChatId(),
			MessageID:   c.Session.LastBotId,
			ReplyMarkup: c.makeInlineKeyboard(buttons),
		},
		Media: tgbotapi.InputMediaPhoto{
			BaseInputMedia: baseMedia,
		},
	}
	c.sendEdit(msg)
}

func (c *Controller) AnswerCallback(text string) {
	c.answerCallback(tgbotapi.NewCallback(c.update.CallbackQuery.ID, text), false)
}

func (c *Controller) AnswerWithAlert(text string) {
	c.answerCallback(tgbotapi.NewCallbackWithAlert(c.update.CallbackQuery.ID, text), true)
}

func (c *Controller) answerCallback(msg tgbotapi.Chattable, isAlert bool) {
	if isAlert {
		c.send(msg)
	} else {
		c.bot.Request(msg)
	}
}

func (c *Controller) ChatId() int64 {
	switch {
	case c.update.Message != nil:
		return c.update.Message.Chat.ID
	case c.update.CallbackQuery != nil:
		return c.update.CallbackQuery.Message.Chat.ID
	case c.update.EditedMessage != nil:
		return c.update.EditedMessage.Chat.ID
	default:
		return 0
	}
}

func (c *Controller) GetUserInfo() tgbotapi.User {
	var user tgbotapi.User
	switch {
	case c.update.Message != nil:
		user = *c.update.Message.From
	case c.update.CallbackQuery != nil:
		user = *c.update.CallbackQuery.From
	case c.update.EditedMessage != nil:
		user = *c.update.EditedMessage.From
	}

	return user
}

func (c *Controller) SendPhotoFile(path, caption string) {
	c.sendPhoto(path, caption, nil)
}

func (c *Controller) SendPhotoFileWithUrl(path, caption string, buttons [][]Button) {
	c.sendPhoto(path, caption, buttons)
}

func (c *Controller) sendPhoto(path, caption string, buttons [][]Button) {
	buff, err := os.ReadFile(path)
	if err != nil {
		log.Println("file read error:", err)
		return
	}
	c.SendPhotoBytes(buff, caption, buttons)
}

func (c *Controller) SendPhotoBytes(buff []byte, caption string, buttons [][]Button) {
	file := tgbotapi.FileBytes{Name: "photo", Bytes: buff}
	msg := tgbotapi.NewPhoto(c.ChatId(), file)
	msg.Caption = caption
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = c.makeInlineKeyboard(buttons)
	c.send(msg)
}

func (c *Controller) SendWithUrl(text string, buttons [][]Button) {
	msg := tgbotapi.NewMessage(c.ChatId(), text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = c.makeInlineKeyboard(buttons)
	c.send(msg)
}

func (c *Controller) send(msg tgbotapi.Chattable) {
	res, err := c.bot.Send(msg)
	if err != nil {
		log.Println("Send error:", err)
		return
	}

	if res.MessageID != 0 {
		c.Session.LastBotId = res.MessageID
	}
}

func (c *Controller) sendEdit(msg tgbotapi.Chattable) {
	_, err := c.bot.Request(msg)
	if err != nil {
		log.Println("SendEdit error:", err)
	}
}

func (c *Controller) makeInlineKeyboard(buttons [][]Button) *tgbotapi.InlineKeyboardMarkup {
	if len(buttons) == 0 {
		return nil
	}
	keyboard := make([][]tgbotapi.InlineKeyboardButton, len(buttons))
	for i, row := range buttons {
		keyboard[i] = make([]tgbotapi.InlineKeyboardButton, len(row))
		for j, btn := range row {
			content := btn.Data
			switch {
			case strings.HasPrefix(btn.Data, "http"):
				keyboard[i][j] = tgbotapi.InlineKeyboardButton{
					Text: btn.Label,
					URL:  &content,
				}
			case strings.HasPrefix(btn.Data, "app:"):
				keyboard[i][j] = tgbotapi.InlineKeyboardButton{
					Text: btn.Label,
					WebApp: &tgbotapi.WebAppInfo{
						URL: strings.Replace(btn.Data, "app", "https", 1),
					},
				}
			default:
				keyboard[i][j] = tgbotapi.InlineKeyboardButton{
					Text:         btn.Label,
					CallbackData: &content,
				}
			}
		}
	}

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
}

func (c *Controller) SendError(err error) {
	c.SendMsg(fmt.Sprintf("❌❌❌ %s 请稍后重试", err.Error()))
}

func (c *Controller) SendInputError(reason string) {
	c.SendMsg(fmt.Sprintf("❌❌❌ %s 请重新输入", reason))
}

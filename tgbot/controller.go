package tgbot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const ParseMode = "HTML"

type Controller struct {
	Session *Session
	Param   string
	Start   string
	User    tgbotapi.User
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
	c.User = c.GetUserInfo()
	c.Session = SessionMgr.GetSession(c.GetUserInfo().ID)

	switch {
	case update.Message != nil:
		c.Session.SaveUserID(update.Message.MessageID)
	case update.CallbackQuery != nil:
		c.Session.SaveUserID(update.CallbackQuery.Message.MessageID)
	case update.EditedMessage != nil:
		c.Session.SaveUserID(update.EditedMessage.MessageID)
	}

	for _, id := range c.Session.GetErrors() {
		c.deleteMessage(id)
	}
	c.Session.ClearErrors()
}

func (c *Controller) Handle() {
	msg := "bot handler not implemented"
	log.Println(msg)
	c.SendError(msg)
}

func (c *Controller) HandleEdit() {
	msg := "bot edit handler not implemented"
	log.Println(msg)
	c.SendError(msg)
}

func (c *Controller) HandleNext() bool {
	msg := "bot anyhandler not implemented"
	log.Println(msg)
	c.SendError(msg)
	return false
}

func (c *Controller) DefaultMenu() {
	c.setMenuButtonUrl("default", "", "")
}

func (c *Controller) ShowMenuUrl(text, url string) {
	c.setMenuButtonUrl("web_app", text, url)
}

func (c *Controller) ShowMenu(text string) {
	c.setMenuButton("commands", text)
}

func (c *Controller) setMenuButton(buttonType, text string) {
	_, err := c.bot.Request(tgbotapi.SetChatMenuButtonConfig{
		ChatID: c.ChatId(),
		MenuButton: &tgbotapi.MenuButton{
			Type: buttonType,
			Text: text,
		},
	})
	if err != nil {
		log.Println("SetMenuButton error:", err)
	}
}

func (c *Controller) setMenuButtonUrl(buttonType, text, url string) {
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
	msg.ParseMode = ParseMode
	msg.ReplyToMessageID = c.update.Message.MessageID
	c.send(msg)
}

func (c *Controller) SendMsg(text string) {
	msg := tgbotapi.NewMessage(c.ChatId(), text)
	msg.ParseMode = ParseMode
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
	c.sendWithoutRecord(msg)
}

func (c *Controller) EditLastBotMsg(text string) {
	msg := tgbotapi.NewEditMessageText(c.ChatId(), c.Session.LastBotId, text)
	msg.ParseMode = ParseMode
	c.sendWithoutRecord(msg)
}

func (c *Controller) EditLastBotMsgWithUrl(text string, buttons [][]Button) {
	msg := tgbotapi.NewEditMessageText(c.ChatId(), c.Session.LastBotId, text)
	msg.ReplyMarkup = c.makeInlineKeyboard(buttons)
	msg.ParseMode = ParseMode
	c.sendWithoutRecord(msg)
}

func (c *Controller) EditLastBotPhoto(path, caption string) {
	c.editMediaWithButtons(path, "photo", caption, nil)
}

func (c *Controller) EditLastBotVideo(path, kind, caption string) {
	c.editMediaWithButtons(path, "video", caption, nil)
}

func (c *Controller) EditLastBotPhotoWithUrl(path, caption string, buttons [][]Button) {
	c.editMediaWithButtons(path, "photo", caption, buttons)
}

func (c *Controller) EditLastBotVideoWithUrl(path, caption string, buttons [][]Button) {
	c.editMediaWithButtons(path, "video", caption, buttons)
}

func (c *Controller) editMediaWithButtons(path, kind, caption string, buttons [][]Button) {
	var file tgbotapi.RequestFileData
	if strings.HasPrefix(path, "http") {
		file = tgbotapi.FileURL(path)
	} else {
		file = tgbotapi.FilePath(path)
	}

	baseMedia := tgbotapi.BaseInputMedia{
		Type:      kind,
		Media:     file,
		Caption:   caption,
		ParseMode: ParseMode,
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
	c.sendWithoutRecord(msg)
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

func (c *Controller) SendPhoto(path, caption string) {
	c.sendPhoto(path, caption, nil, false)
}

func (c *Controller) SendPhotoFileWithUrl(path, caption string, buttons [][]Button) {
	c.sendPhoto(path, caption, buttons, false)
}

func (c *Controller) SendPhotoFileWithKeyboard(path, caption string, buttons [][]Button) {
	c.sendPhoto(path, caption, buttons, true)
}

func (c *Controller) sendPhoto(path, caption string, buttons [][]Button, keyboard bool) {
	var msg tgbotapi.PhotoConfig
	log.Println("path:", path)
	if strings.HasPrefix(path, "http") {
		log.Println("path:", path)
		msg = tgbotapi.NewPhoto(c.ChatId(), tgbotapi.FileURL(path))
	} else {
		msg = tgbotapi.NewPhoto(c.ChatId(), tgbotapi.FilePath(path))
	}

	msg.Caption = caption
	msg.ParseMode = ParseMode
	if keyboard {
		msg.ReplyMarkup = c.makeKeyboard(buttons)
		c.sendWithoutRecord(msg)
	} else {
		msg.ReplyMarkup = c.makeInlineKeyboard(buttons)
		c.send(msg)
	}
}

func (c *Controller) SendVideo(path, caption string) {
	c.sendVideo(path, caption, nil, false)
}

func (c *Controller) SendVideoWithUrl(path, caption string, buttons [][]Button) {
	c.sendVideo(path, caption, buttons, false)
}

func (c *Controller) SendVideoWithKeyboard(path, caption string, buttons [][]Button) {
	c.sendVideo(path, caption, buttons, true)
}

func (c *Controller) sendVideo(path, caption string, buttons [][]Button, keyboard bool) {
	msg := tgbotapi.NewVideo(c.ChatId(), tgbotapi.FilePath(path))
	msg.Caption = caption
	msg.ParseMode = ParseMode
	if keyboard {
		msg.ReplyMarkup = c.makeKeyboard(buttons)
	} else {
		msg.ReplyMarkup = c.makeInlineKeyboard(buttons)
	}
	c.send(msg)
}

func (c *Controller) SendWithUrl(text string, buttons [][]Button) {
	c.sendMsg(text, buttons, false)
}

func (c *Controller) SendWithKeyboard(text string, buttons [][]Button) {
	c.sendMsg(text, buttons, true)
}

func (c *Controller) sendMsg(text string, buttons [][]Button, keyboard bool) {
	msg := tgbotapi.NewMessage(c.ChatId(), text)
	msg.ParseMode = ParseMode
	if keyboard {
		msg.ReplyMarkup = c.makeKeyboard(buttons)
		c.send(msg)
	} else {
		msg.ReplyMarkup = c.makeInlineKeyboard(buttons)
		c.send(msg)
	}
}

func (c *Controller) send(msg tgbotapi.Chattable) {
	res, err := c.bot.Send(msg)
	if err != nil {
		log.Println("Send error:", err)
		return
	}

	log.Println("res:", res.MessageID)
	if res.MessageID != 0 {
		c.Session.LastBotId = res.MessageID
	}
}

func (c *Controller) sendError(text string) {
	msg := tgbotapi.NewMessage(c.ChatId(), text)
	res, err := c.bot.Send(msg)
	if err != nil {
		log.Println("Send error:", err)
		return
	}

	if res.MessageID != 0 {
		c.Session.AddError(res.MessageID)
	}
}

func (c *Controller) sendWithoutRecord(msg tgbotapi.Chattable) {
	_, err := c.bot.Request(msg)
	if err != nil {
		log.Println("SendWithoutRecord error:", err)
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
			var button tgbotapi.InlineKeyboardButton
			switch {
			case strings.HasPrefix(btn.Data, "http"):
				button = tgbotapi.InlineKeyboardButton{
					Text: btn.Label,
					URL:  &content,
				}
			case strings.HasPrefix(btn.Data, "app:"):
				button = tgbotapi.InlineKeyboardButton{
					Text: btn.Label,
					WebApp: &tgbotapi.WebAppInfo{
						URL: strings.Replace(btn.Data, "app", "https", 1),
					},
				}
			default:
				button = tgbotapi.InlineKeyboardButton{
					Text:         btn.Label,
					CallbackData: &content,
				}
			}
			keyboard[i][j] = button
		}
	}

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
}

func (c *Controller) makeKeyboard(buttons [][]Button) *tgbotapi.ReplyKeyboardMarkup {
	keyboard := make([][]tgbotapi.KeyboardButton, len(buttons))
	for i, row := range buttons {
		keyboard[i] = make([]tgbotapi.KeyboardButton, len(row))
		for j, btn := range row {
			var button tgbotapi.KeyboardButton
			if strings.HasPrefix(btn.Data, "app") {
				button = tgbotapi.KeyboardButton{
					Text: btn.Label,
					WebApp: &tgbotapi.WebAppInfo{
						URL: strings.Replace(btn.Data, "app", "https", 1),
					},
				}
			} else {
				if btn.Data == "phone" {
					button = tgbotapi.NewKeyboardButtonContact(btn.Label)
				} else if btn.Data == "geo" {
					button = tgbotapi.NewKeyboardButtonLocation(btn.Label)
				} else {
					button = tgbotapi.KeyboardButton{
						Text: btn.Label,
					}
				}
			}
			keyboard[i][j] = button
		}
	}

	return &tgbotapi.ReplyKeyboardMarkup{
		Keyboard:       keyboard,
		ResizeKeyboard: true,
	}
}

func (c *Controller) RemoveKeyboard(text string) {
	msg := tgbotapi.NewMessage(c.ChatId(), text)
	msg.ParseMode = ParseMode
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	c.send(msg)
}

func (c *Controller) SendError(reason string) {
	c.sendError(fmt.Sprintf("❌❌❌ %s 请稍后重试", reason))
}

func (c *Controller) SendInputError(reason string) {
	c.sendError(fmt.Sprintf("❌❌❌ %s 请重新输入", reason))
}

func (c *Controller) Contact() *tgbotapi.Contact {
	if c.update.Message == nil {
		return nil
	}

	return c.update.Message.Contact
}

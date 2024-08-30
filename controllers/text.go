package controllers

import "github.com/scnon/tg-bot/tgbot"

// 文本消息处理器
type TextController struct {
	tgbot.Controller
}

// 处理文本消息
func (c *TextController) Handle() {
	c.RemoveKeyboard("remove")
}

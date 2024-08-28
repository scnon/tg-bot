package controllers

import "tg-bot/tgbot"

// 文本消息处理器
type TextController struct {
	tgbot.Controller
}

// 处理文本消息
func (c *TextController) Handle() {
	c.SendMsg("收到文本消息")
}

package controllers

import (
	"github.com/scnon/tg-bot/tgbot"
)

// 文本消息处理器
type TextController struct {
	tgbot.Controller
}

// 处理文本消息
func (c *TextController) Handle() {
	c.DeleteLastUserMsg()
	// c.EditLastBotMsgWithUrl("remove", [][]tgbot.Button{
	// 	{
	// 		{Id: 1, Label: "删除", Data: "remove"},
	// 	},
	// })
	if c.Param != "1" {
		c.SendInputError("参数错误")
	} else {
		c.SendMsg("start controller handle")
	}
}

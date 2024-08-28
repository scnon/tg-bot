package controllers

import (
	"log"
	"tg-bot/tgbot"
)

// start 命令处理器
type StartController struct {
	tgbot.Controller
}

// 处理 start 命令
func (c *StartController) Handle() {
	log.Println("start controller handle", c.Param)
	c.SendWithUrl("欢迎使用机器人", [][]tgbot.Button{
		{
			tgbot.Button{Label: "查询", Data: "query"},
		},
	})
}

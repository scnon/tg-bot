package controllers

import (
	"log"

	"github.com/scnon/tg-bot/tgbot"
)

// start 命令处理器
type StartController struct {
	tgbot.Controller
}

// 处理 start 命令
func (c *StartController) Handle() {
	log.Println("start controller handle", c.Param)
	c.SendWithUrl("calc", [][]tgbot.Button{
		{
			{Id: 1, Label: "加法", Data: "add"},
		},
	})

}

package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/scnon/tg-bot/tgbot"
)

type AddController struct {
	tgbot.Controller
}

func (c *AddController) Handle() {
	params := strings.Split(c.Param, " ")
	result := calcResult(params)

	c.EditLastBotPhotoWithUrl("assets/edit_1.png", fmt.Sprintf("结果是：%d", result), [][]tgbot.Button{
		{
			{Id: 1, Label: "编辑", Data: "edit"},
		},
	})
}

func (c *AddController) HandleEdit() {
	params := strings.Split(c.Param, " ")
	result := calcResult(params)

	c.EditLastBotMsg(fmt.Sprintf("结果是：%d", result))
}

func calcResult(params []string) int64 {
	var result int64
	for _, v := range params {
		num, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			result += num
		}
	}
	return result
}

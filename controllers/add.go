package controllers

import (
	"fmt"
	"github.com/scnon/tg-bot/tgbot"
	"strconv"
	"strings"
)

type AddController struct {
	tgbot.Controller
}

func (c *AddController) Handle() {
	params := strings.Split(c.Param, " ")
	result := calcResult(params)

	c.SendMsg(fmt.Sprintf("结果是：%d", result))
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

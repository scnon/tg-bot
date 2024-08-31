package controllers

import (
	"errors"

	"github.com/scnon/tg-bot/tgbot"
)

// 查询处理器
type QueryController struct {
	tgbot.Controller
}

// 处理查询消息
func (c *QueryController) Handle() {
	// c.SendMsg("请输入查询内容")
	c.SendError(errors.New("输入错误，请重新输入"))
}

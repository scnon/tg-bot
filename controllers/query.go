package controllers

import (
	"fmt"

	"github.com/scnon/tg-bot/tgbot"
)

// 查询处理器
type QueryController struct {
	tgbot.Controller
}

// 处理查询消息
func (c *QueryController) Handle() {
	// c.SendMsg("请输入查询内容")
	// c.SendInputError("输入错误")

	c.EditLastBotVideoWithUrl("assets/business.mp4", fmt.Sprintf("结果是：%d", 3), [][]tgbot.Button{
		{
			{Id: 1, Label: "编辑", Data: "edit"},
		},
	})
}

package main

import (
	"log"

	"github.com/scnon/tg-bot/controllers"
	"github.com/scnon/tg-bot/tgbot"
)

func main() {
	token := "6588358108:AAGgEe8wDgcHR_f8iSi90GU-uS-u1W47wjE" // os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Println("BOT_TOKEN", "not found")
		return
	}

	bot, err := tgbot.NewTgBot(token)
	if err != nil {
		log.Println("NewTgBot error:", err)
		return
	}

	bot.RegCmdController("start", &controllers.StartController{})
	bot.RegCmdController("add", &controllers.AddController{})
	bot.RegQueryController(&controllers.QueryController{})
	bot.RegTextController(&controllers.TextController{})

	bot.StartTgLoop()
}

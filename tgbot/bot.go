package tgbot

import (
	"encoding/json"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Button struct {
	Id    int32  `json:"id" comment:"id"`
	Label string `json:"label" comment:"按钮名称"`
	Data  string `json:"data" comment:"按钮数据"`
}

type TgBot struct {
	token         string
	bot           *tgbotapi.BotAPI
	anyController ControllerInterface
	cmds          map[string]ControllerInterface
	texts         []ControllerInterface
	querys        []ControllerInterface
}

func NewTgBot(token string) (*TgBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println("Telegram bot init error:", err)
		return nil, err
	}

	log.Printf("Authorized on account [%s] @%s", bot.Self.FirstName, bot.Self.UserName)
	return &TgBot{
		token:  token,
		bot:    bot,
		cmds:   make(map[string]ControllerInterface),
		querys: []ControllerInterface{},
		texts:  []ControllerInterface{},
	}, nil
}

func (b *TgBot) RegAnyController(controller ControllerInterface) {
	b.anyController = controller
}

func (b *TgBot) RegCmdController(cmd string, controller ControllerInterface) {
	b.cmds[cmd] = controller
}

func (b *TgBot) RegTextController(controller ControllerInterface) {
	b.texts = append(b.texts, controller)
}

func (b *TgBot) RegQueryController(controller ControllerInterface) {
	b.querys = append(b.querys, controller)
}

func (b *TgBot) RemoveController(controller ControllerInterface) {
	remove := func(controllers []ControllerInterface) []ControllerInterface {
		for i, c := range controllers {
			if c == controller {
				return append(controllers[:i], controllers[i+1:]...)
			}
		}
		return controllers
	}

	b.texts = remove(b.texts)
	b.querys = remove(b.querys)
}

func (b *TgBot) StartTgLoop() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if b.anyController != nil {
			b.anyController.Init(update, b.bot, "")
			if !b.anyController.HandleNext() {
				continue
			}
		}

		switch {
		case update.CallbackQuery != nil:
			query := update.CallbackQuery.Data
			for _, handler := range b.querys {
				handler.Init(update, b.bot, query)
				go handler.Handle()
			}
		case update.Message != nil,
			update.EditedMessage != nil:
			isEdit := update.EditedMessage != nil
			message := update.Message
			if isEdit {
				message = update.EditedMessage
			}

			text := message.Text
			if message.IsCommand() {
				command := message.Command()
				if handler, ok := b.cmds[command]; ok {
					handler.Init(update, b.bot, text)
					if isEdit {
						go handler.HandleEdit()
					} else {
						go handler.Handle()
					}
				}
			} else {
				for _, handler := range b.texts {
					handler.Init(update, b.bot, text)
					if isEdit {
						go handler.HandleEdit()
					} else {
						go handler.Handle()
					}
				}
			}
		default:
			res, err := json.Marshal(update)
			if err != nil {
				log.Println("json marshal error:", err)
				return
			}
			log.Println("Unhandled update:", string(res))
		}
	}
}

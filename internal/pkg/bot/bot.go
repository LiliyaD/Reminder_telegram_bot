package bot

import (
	"context"
	"fmt"

	configPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/config"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	commandPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/bot/command"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Interface interface {
	InitCommands()
	Run() error
}

func MustNew(client pb.AdminClient) Interface {
	bot, err := tgBotAPI.NewBotAPI(configPkg.GetApiKey())
	if err != nil {
		journal.LogFatal(errors.Wrap(err, "init tgbot"))
	}

	bot.Debug = true

	return &commander{
		bot:    bot,
		route:  make(map[string]commandPkg.InterfaceCommand),
		client: client,
	}
}

type commander struct {
	bot    *tgBotAPI.BotAPI
	route  map[string]commandPkg.InterfaceCommand
	client pb.AdminClient
}

func (c *commander) InitCommands() {
	c.route = commandPkg.CreateCommands()
}

func (c *commander) Run() error {
	u := tgBotAPI.NewUpdate(0)
	updates := c.bot.GetUpdatesChan(u)

	ctx := context.Background()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgBotAPI.NewMessage(update.Message.Chat.ID, "")

		if cmd := update.Message.Command(); cmd != "" {
			if cmdName, ok := c.route[cmd]; ok {
				chat := models.Chat{
					ChatID:   update.Message.Chat.ID,
					UserName: update.Message.Chat.FirstName,
				}
				msg.Text = cmdName.Process(ctx, update.Message.CommandArguments(), chat, c.client)
			} else {
				msg.Text = "Unknown command"
			}
		} else {
			msg.Text = fmt.Sprintf("you send <%v>", update.Message.Text)
		}
		_, err := c.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "send tg message")
		}
	}
	return nil
}

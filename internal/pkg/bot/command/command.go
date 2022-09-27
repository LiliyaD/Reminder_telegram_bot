package command

import (
	"context"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
)

type InterfaceCommand interface {
	Process(ctx context.Context, s string, chat models.Chat, client pb.AdminClient) string
	Name() string
}

var (
	route map[string]InterfaceCommand
)

func CreateCommands() map[string]InterfaceCommand {
	route = make(map[string]InterfaceCommand)
	addHandler(&commandStart{})
	addHandler(&commandHelp{})
	addHandler(&commandExamples{})
	//addHandler(&commandRemind{})
	addHandler(&commandList{})
	addHandler(&commandStream{})
	addHandler(&commandGet{})
	addHandler(&commandToday{})
	addHandler(&commandCreate{})
	addHandler(&commandUpdate{})
	addHandler(&commandDelete{})
	return route
}

func addHandler(cmd InterfaceCommand) {
	route[cmd.Name()] = cmd
}

package command

import (
	"context"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
)

type commandHelp struct {
	InterfaceCommand
}

func (c *commandHelp) Process(_ context.Context, _ string, _ models.Chat, _ pb.AdminClient) string {
	return "Available commands: \n" +
		"/help - list commands\n" +
		"/examples - show examples\n" +
		//"/remind <time> - bot will remind today's activities at the exact time\n" +
		"/list, /list_stream - list all activities\n" +
		"/get <name> - get information about activity\n" +
		"/today - list today activities\n" +
		"/add <name> <begin_date> <end_date> <times_per_day> <quantity_per_time> - add new activity\n" +
		"/delete <name> - delete activity\n" +
		"/update <name> <begin_date> <end_date> <times_per_day> <quantity_per_time> - update activity, if you don't want to change all fields write _ for parameters, which you won't change"
}

func (c *commandHelp) Name() string {
	return "help"
}

package command

import (
	"context"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
)

type commandStart struct {
	InterfaceCommand
}

func (c *commandStart) Process(_ context.Context, _ string, _ models.Chat, _ pb.AdminClient) string {
	return "Hi! This bot will remind you some things, that you should do each day.\n" +
		"Write /help to see available commands."
}

func (c *commandStart) Name() string {
	return "start"
}

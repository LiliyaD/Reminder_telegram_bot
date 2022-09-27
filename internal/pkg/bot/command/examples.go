package command

import (
	"context"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
)

type commandExamples struct {
	InterfaceCommand
}

func (c *commandExamples) Process(_ context.Context, _ string, _ models.Chat, _ pb.AdminClient) string {
	return "Example 1: If you want to add reminder that you want to do 30 push-ups 2 times per day from 20.07.2022 to 30.07.2022 write in chat:\n" +
		"/add push-up 20.07.2022 30.07.2022 2 30\n\n" +
		"Example 2: If you want to add reminder that you have to give medicine to your pet 2 times per day in dosage 1/4 of pill from 20.07.2022 to 30.07.2022 write in chat:\n" +
		"/add medicine 20.07.2022 30.07.2022 2 0.5\n"
}

func (c *commandExamples) Name() string {
	return "examples"
}

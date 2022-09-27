package command

import (
	"context"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	commonPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/common/commands"
	multiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/multithread"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type commandToday struct {
	InterfaceCommand
}

func (c *commandToday) Process(ctx context.Context, s string, chat models.Chat, client pb.AdminClient) string {
	response, err := commonPkg.Today(ctx, chat.ChatID, client)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			journal.LogError(multiPkg.ErrNoActivitiesForToday)
			return multiPkg.ErrNoActivitiesForToday.Error()
		}
		journal.LogError(err)
		return prepareErrorTextForUser(err)
	}

	var lines string
	for _, v := range response.Activities {
		lines += formText(v.GetName(), v.GetBeginDate(), v.GetEndDate(), v.GetTimesPerDay(), v.GetQuantityPerTime())
	}

	if len(lines) == 0 {
		journal.LogError(multiPkg.ErrNoActivitiesForToday)
		return multiPkg.ErrNoActivitiesForToday.Error()
	}

	return lines
}

func (c *commandToday) Name() string {
	return "today"
}

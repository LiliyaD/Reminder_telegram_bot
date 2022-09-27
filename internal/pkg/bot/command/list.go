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

type commandList struct {
	InterfaceCommand
}

func (c *commandList) Process(ctx context.Context, s string, chat models.Chat, client pb.AdminClient) string {
	response, err := commonPkg.List(ctx, chat.ChatID, models.Pagination{}, client)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			journal.LogError(multiPkg.ErrEmptyStorage)
			return multiPkg.ErrEmptyStorage.Error()
		}
		journal.LogError(err)
		return prepareErrorTextForUser(err)
	}

	var lines string
	for _, v := range response.Activities {
		lines += formText(v.GetName(), v.GetBeginDate(), v.GetEndDate(), v.GetTimesPerDay(), v.GetQuantityPerTime())
	}

	if len(lines) == 0 {
		journal.LogError(multiPkg.ErrEmptyStorage)
		return multiPkg.ErrEmptyStorage.Error()
	}

	return lines
}

func (c *commandList) Name() string {
	return "list"
}

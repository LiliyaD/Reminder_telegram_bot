package command

import (
	"context"
	"fmt"
	"time"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	commonPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/common/commands"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
)

type commandCreate struct {
	InterfaceCommand
}

func (c *commandCreate) Process(ctx context.Context, s string, chat models.Chat, client pb.AdminClient) string {
	res, err := parseData(s, false)
	if err != nil {
		journal.LogError(err)
		return err.Error()
	}

	response, err := commonPkg.Create(ctx, res[0].(string),
		models.DailyActivity{
			BeginDate:       res[1].(time.Time),
			EndDate:         res[2].(time.Time),
			TimesPerDay:     res[3].(uint8),
			QuantityPerTime: res[4].(float32)}, chat, client)
	if err != nil {
		journal.LogError(err)
		return prepareErrorTextForUser(err)
	}

	if response == nil {
		return "Request is accepted for processing"
	}

	return fmt.Sprintf("Activity %v was added:\n begin date %v \n end date %v \n times per day %v \n quantity per time %v",
		response.GetName(),
		response.GetBeginDate(),
		response.GetEndDate(),
		response.GetTimesPerDay(),
		response.GetQuantityPerTime())
}

func (c *commandCreate) Name() string {
	return "add"
}

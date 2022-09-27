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

type commandUpdate struct {
	InterfaceCommand
}

func (c *commandUpdate) Process(ctx context.Context, s string, chat models.Chat, client pb.AdminClient) string {
	res, err := parseData(s, true)
	if err != nil {
		journal.LogError(err)
		return err.Error()
	}

	response, err := commonPkg.Update(ctx, chat.ChatID, res[0].(string),
		models.DailyActivity{
			BeginDate:       res[1].(time.Time),
			EndDate:         res[2].(time.Time),
			TimesPerDay:     res[3].(uint8),
			QuantityPerTime: res[4].(float32)}, client)
	if err != nil {
		journal.LogError(err)
		return prepareErrorTextForUser(err)
	}

	if response == nil {
		return "Request is accepted for processing"
	}

	return fmt.Sprintf("Activity %v was updated:\n begin date %v \n end date %v \n times per day %v \n quantity per time %v",
		response.GetName(),
		response.GetBeginDate(),
		response.GetEndDate(),
		response.GetTimesPerDay(),
		response.GetQuantityPerTime())
}

func (c *commandUpdate) Name() string {
	return "update"
}

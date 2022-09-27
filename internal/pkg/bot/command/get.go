package command

import (
	"context"
	"strings"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	commonPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/common/commands"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
)

type commandGet struct {
	InterfaceCommand
}

func (c *commandGet) Process(ctx context.Context, s string, chat models.Chat, client pb.AdminClient) string {
	if len(s) == 0 {
		journal.LogError(errBadParameters)
		return errBadParameters.Error()
	}

	p := strings.Split(s, " ")
	if len(p) != 1 {
		journal.LogError(errBadParameters)
		return errBadParameters.Error()
	}

	v, err := commonPkg.Get(ctx, chat.ChatID, p[0], client)
	if err != nil {
		journal.LogError(err)
		return prepareErrorTextForUser(err)
	}

	return formText(v.GetName(), v.GetBeginDate(), v.GetEndDate(), v.GetTimesPerDay(), v.GetQuantityPerTime())
}

func (c *commandGet) Name() string {
	return "get"
}

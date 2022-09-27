package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	commonPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/common/commands"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
)

type commandDelete struct {
	InterfaceCommand
}

func (c *commandDelete) Process(ctx context.Context, s string, chat models.Chat, client pb.AdminClient) string {
	if len(s) == 0 {
		journal.LogError(errBadParameters)
		return errBadParameters.Error()
	}

	p := strings.Split(s, " ")
	if len(p) != 1 {
		journal.LogError(errBadParameters)
		return errBadParameters.Error()
	}

	err := commonPkg.Delete(ctx, chat.ChatID, p[0], client)
	if err != nil {
		journal.LogError(err)
		return prepareErrorTextForUser(err)
	}

	return fmt.Sprintf("Activity %v was deleted", p[0])
}

func (c *commandDelete) Name() string {
	return "delete"
}

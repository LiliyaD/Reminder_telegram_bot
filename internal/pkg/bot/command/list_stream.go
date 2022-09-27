// ---------------------------------------
//
// Special variant for gRPC stream
//
// ---------------------------------------

package command

import (
	"context"
	"io"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	commonPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/common/commands"
	multiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/multithread"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type commandStream struct {
	InterfaceCommand
}

func (c *commandStream) Process(ctx context.Context, s string, chat models.Chat, client pb.AdminClient) string {
	stream, err := commonPkg.ListStream(ctx, chat.ChatID, client)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			journal.LogError(multiPkg.ErrEmptyStorage)
			return multiPkg.ErrEmptyStorage.Error()
		}
		journal.LogError(err)
		return prepareErrorTextForUser(err)
	}

	var lines string
	for {
		v, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			if status.Code(err) == codes.NotFound {
				journal.LogError(multiPkg.ErrEmptyStorage)
				return multiPkg.ErrEmptyStorage.Error()
			}
			journal.LogError(err)
			return prepareErrorTextForUser(err)
		}
		lines += formText(v.GetName(), v.GetBeginDate(), v.GetEndDate(), v.GetTimesPerDay(), v.GetQuantityPerTime())
	}
	return lines
}

func (c *commandStream) Name() string {
	return "list_stream"
}

package cache

import (
	"context"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
)

type Interface interface {
	Add(ctx context.Context, name string, act models.DailyActivity, chat models.Chat) (models.DailyActivity, error)
	Delete(ctx context.Context, name string, chatID int64) error
	Update(ctx context.Context, name string, act models.DailyActivity, chatID int64) (models.DailyActivity, error)
	List(ctx context.Context, chatID int64, pagination models.Pagination) (map[string]models.DailyActivity, error)
	ListStream(ctx context.Context, chatID int64, ch chan models.DailyActivityRec, chErr chan error)
	Read(ctx context.Context, name string, chatID int64) (models.DailyActivity, error)
	ReadToday(ctx context.Context, chatID int64) (map[string]models.DailyActivity, error)
}

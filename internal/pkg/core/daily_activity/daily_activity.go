package daily_activity

import (
	"context"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	cachePkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache"
	multiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/multithread"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
)

type Interface interface {
	Add(ctx context.Context, name string, act models.DailyActivity, chat models.Chat) (models.DailyActivity, error)
	Delete(ctx context.Context, name string, chatID int64) error
	Update(ctx context.Context, name string, act models.DailyActivity, chatID int64) (models.DailyActivity, error)
	List(ctx context.Context, chatID int64, pagination models.Pagination) (map[string]models.DailyActivity, error)
	ListStream(ctx context.Context, chatID int64, ch chan models.DailyActivityRec, chErr chan error)
	Get(ctx context.Context, name string, chatID int64) (models.DailyActivity, error)
	Today(ctx context.Context, chatID int64) (map[string]models.DailyActivity, error)
}

func New(dataStorage cachePkg.Interface) Interface {
	return &core{
		cache: multiPkg.NewTimeoutStorage(multiPkg.NewConcurrencyStorage(dataStorage)),
	}
}

type core struct {
	cache cachePkg.Interface
}

func (c *core) Add(ctx context.Context, name string, act models.DailyActivity, chat models.Chat) (models.DailyActivity, error) {
	act, err := c.cache.Add(ctx, name, act, chat)
	if err != nil {
		journal.LogError(err)
	}
	return act, err
}

func (c *core) Delete(ctx context.Context, name string, chatID int64) error {
	err := c.cache.Delete(ctx, name, chatID)
	if err != nil {
		journal.LogError(err)
	}
	return err
}

func (c *core) Update(ctx context.Context, name string, act models.DailyActivity, chatID int64) (models.DailyActivity, error) {
	act, err := c.cache.Update(ctx, name, act, chatID)
	if err != nil {
		journal.LogError(err)
	}
	return act, err
}

func (c *core) List(ctx context.Context, chatID int64, pagination models.Pagination) (map[string]models.DailyActivity, error) {
	act, err := c.cache.List(ctx, chatID, pagination)
	if err != nil {
		journal.LogError(err)
	}
	return act, err
}

func (c *core) ListStream(ctx context.Context, chatID int64, ch chan models.DailyActivityRec, chErr chan error) {
	c.cache.ListStream(ctx, chatID, ch, chErr)
}

func (c *core) Get(ctx context.Context, name string, chatID int64) (models.DailyActivity, error) {
	act, err := c.cache.Read(ctx, name, chatID)
	if err != nil {
		journal.LogError(err)
	}
	return act, err
}

func (c *core) Today(ctx context.Context, chatID int64) (map[string]models.DailyActivity, error) {
	act, err := c.cache.ReadToday(ctx, chatID)
	if err != nil {
		journal.LogError(err)
	}
	return act, err
}

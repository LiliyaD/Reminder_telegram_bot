package multithread

import (
	"context"
	"time"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	"github.com/pkg/errors"
)

const taskTimeout = 2 * time.Second

var ErrTimeOut = errors.New("Timeout error")

type TimeoutStorage struct {
	storage InterfaceStorage
}

func NewTimeoutStorage(storage InterfaceStorage) InterfaceStorage {
	return &TimeoutStorage{
		storage: storage,
	}
}

func (c *TimeoutStorage) List(ctx context.Context, chatID int64, pagination models.Pagination) (map[string]models.DailyActivity, error) {
	ctx, cancel := context.WithTimeout(ctx, taskTimeout) //Время ожидания
	defer cancel()
	act, err := c.storage.List(ctx, chatID, pagination)
	if err != nil {
		journal.LogError(err)
	}
	return act, err
}

func (c *TimeoutStorage) ListStream(ctx context.Context, chatID int64, ch chan models.DailyActivityRec, chErr chan error) {
	ctx, cancel := context.WithTimeout(ctx, taskTimeout) //Время ожидания
	defer cancel()
	c.storage.ListStream(ctx, chatID, ch, chErr)
}

func (c *TimeoutStorage) Add(ctx context.Context, name string, act models.DailyActivity, chat models.Chat) (models.DailyActivity, error) {
	ctx, cancel := context.WithTimeout(ctx, taskTimeout) //Время ожидания
	defer cancel()
	act, err := c.storage.Add(ctx, name, act, chat)
	if err != nil {
		journal.LogError(err)
	}
	return act, err
}

func (c *TimeoutStorage) Update(ctx context.Context, name string, newAct models.DailyActivity, chatID int64) (models.DailyActivity, error) {
	ctx, cancel := context.WithTimeout(ctx, taskTimeout) //Время ожидания
	defer cancel()
	act, err := c.storage.Update(ctx, name, newAct, chatID)
	if err != nil {
		journal.LogError(err)
	}
	return act, err
}

func (c *TimeoutStorage) Delete(ctx context.Context, name string, chatID int64) error {
	ctx, cancel := context.WithTimeout(ctx, taskTimeout) //Время ожидания
	defer cancel()
	err := c.storage.Delete(ctx, name, chatID)
	if err != nil {
		journal.LogError(err)
	}
	return err
}

func (c *TimeoutStorage) Read(ctx context.Context, name string, chatID int64) (models.DailyActivity, error) {
	ctx, cancel := context.WithTimeout(ctx, taskTimeout) //Время ожидания
	defer cancel()
	act, err := c.storage.Read(ctx, name, chatID)
	if err != nil {
		journal.LogError(err)
	}
	return act, err
}

func (c *TimeoutStorage) ReadToday(ctx context.Context, chatID int64) (map[string]models.DailyActivity, error) {
	ctx, cancel := context.WithTimeout(ctx, taskTimeout) //Время ожидания
	defer cancel()
	act, err := c.storage.ReadToday(ctx, chatID)
	if err != nil {
		journal.LogError(err)
	}
	return act, err
}

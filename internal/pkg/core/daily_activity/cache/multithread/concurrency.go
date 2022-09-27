package multithread

import (
	"context"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
)

const poolSize = 10

type ConcurrencyStorage struct {
	storage InterfaceStorage
	pool    chan bool
}

func NewConcurrencyStorage(storage InterfaceStorage) InterfaceStorage {
	return &ConcurrencyStorage{
		storage: storage,
		pool:    make(chan bool, poolSize), //Количество соединений (workers)
	}
}

func (c *ConcurrencyStorage) List(ctx context.Context, chatID int64, pagination models.Pagination) (map[string]models.DailyActivity, error) {
	select {
	case c.pool <- true:
		defer func() {
			<-c.pool
		}()
		return c.storage.List(ctx, chatID, pagination)
	case <-ctx.Done():
		// не дождались пула
		journal.LogError(ErrTimeOut)
		return nil, ErrTimeOut
	}
}

func (c *ConcurrencyStorage) ListStream(ctx context.Context, chatID int64, ch chan models.DailyActivityRec, chErr chan error) {
	select {
	case c.pool <- true:
		defer func() {
			<-c.pool
		}()
		c.storage.ListStream(ctx, chatID, ch, chErr)
	case <-ctx.Done():
		// не дождались пула
		journal.LogError(ErrTimeOut)
		chErr <- ErrTimeOut
	}
}

func (c *ConcurrencyStorage) Add(ctx context.Context, name string, act models.DailyActivity, chat models.Chat) (models.DailyActivity, error) {
	select {
	case c.pool <- true:
		defer func() {
			<-c.pool
		}()
		return c.storage.Add(ctx, name, act, chat)
	case <-ctx.Done():
		// не дождались пула
		journal.LogError(ErrTimeOut)
		return models.DailyActivity{}, ErrTimeOut
	}
}

func (c *ConcurrencyStorage) Update(ctx context.Context, name string, newAct models.DailyActivity, chatID int64) (models.DailyActivity, error) {
	select {
	case c.pool <- true:
		defer func() {
			<-c.pool
		}()
		return c.storage.Update(ctx, name, newAct, chatID)
	case <-ctx.Done():
		// не дождались пула
		journal.LogError(ErrTimeOut)
		return models.DailyActivity{}, ErrTimeOut
	}
}

func (c *ConcurrencyStorage) Delete(ctx context.Context, name string, chatID int64) error {
	select {
	case c.pool <- true:
		defer func() {
			<-c.pool
		}()
		return c.storage.Delete(ctx, name, chatID)
	case <-ctx.Done():
		// не дождались пула
		journal.LogError(ErrTimeOut)
		return ErrTimeOut
	}
}

func (c *ConcurrencyStorage) Read(ctx context.Context, name string, chatID int64) (models.DailyActivity, error) {
	select {
	case c.pool <- true:
		defer func() {
			<-c.pool
		}()
		return c.storage.Read(ctx, name, chatID)
	case <-ctx.Done():
		// не дождались пула
		journal.LogError(ErrTimeOut)
		return models.DailyActivity{}, ErrTimeOut
	}
}

func (c *ConcurrencyStorage) ReadToday(ctx context.Context, chatID int64) (map[string]models.DailyActivity, error) {
	select {
	case c.pool <- true:
		defer func() {
			<-c.pool
		}()
		return c.storage.ReadToday(ctx, chatID)
	case <-ctx.Done():
		// не дождались пула
		journal.LogError(ErrTimeOut)
		return nil, ErrTimeOut
	}
}

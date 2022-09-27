package local

import (
	"context"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	multiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/multithread"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
)

type MutexStorage struct {
	storage InterfaceStorage
	mutex   chan bool
}

func NewMutexStorage(storage InterfaceStorage) InterfaceStorage {
	return &MutexStorage{
		storage: storage,
		mutex:   make(chan bool, 1),
	}
}

func (c *MutexStorage) List(ctx context.Context, chatID int64, pagination models.Pagination) (map[string]models.DailyActivity, error) {
	select {
	case c.mutex <- true:
		defer func() {
			<-c.mutex
		}()
		return c.storage.List(ctx, chatID, pagination)
	case <-ctx.Done():
		// не дождались мьютекса
		journal.LogError(multiPkg.ErrTimeOut)
		return nil, multiPkg.ErrTimeOut
	}
}

func (c *MutexStorage) ListStream(ctx context.Context, chatID int64, ch chan models.DailyActivityRec, chErr chan error) {
	select {
	case c.mutex <- true:
		defer func() {
			<-c.mutex
		}()
		c.storage.ListStream(ctx, chatID, ch, chErr)
	case <-ctx.Done():
		// не дождались мьютекса
		journal.LogError(multiPkg.ErrTimeOut)
		chErr <- multiPkg.ErrTimeOut
	}
}

func (c *MutexStorage) Add(ctx context.Context, name string, act models.DailyActivity, chat models.Chat) (models.DailyActivity, error) {
	select {
	case c.mutex <- true:
		defer func() {
			<-c.mutex
		}()
		return c.storage.Add(ctx, name, act, chat)
	case <-ctx.Done():
		// не дождались мьютекса
		journal.LogError(multiPkg.ErrTimeOut)
		return models.DailyActivity{}, multiPkg.ErrTimeOut
	}
}

func (c *MutexStorage) Update(ctx context.Context, name string, newAct models.DailyActivity, chatID int64) (models.DailyActivity, error) {
	select {
	case c.mutex <- true:
		defer func() {
			<-c.mutex
		}()
		return c.storage.Update(ctx, name, newAct, chatID)
	case <-ctx.Done():
		// не дождались мьютекса
		journal.LogError(multiPkg.ErrTimeOut)
		return models.DailyActivity{}, multiPkg.ErrTimeOut
	}
}

func (c *MutexStorage) Delete(ctx context.Context, name string, chatID int64) error {
	select {
	case c.mutex <- true:
		defer func() {
			<-c.mutex
		}()
		return c.storage.Delete(ctx, name, chatID)
	case <-ctx.Done():
		// не дождались мьютекса
		journal.LogError(multiPkg.ErrTimeOut)
		return multiPkg.ErrTimeOut
	}
}

func (c *MutexStorage) Read(ctx context.Context, name string, chatID int64) (models.DailyActivity, error) {
	select {
	case c.mutex <- true:
		defer func() {
			<-c.mutex
		}()
		return c.storage.Read(ctx, name, chatID)
	case <-ctx.Done():
		// не дождались мьютекса
		journal.LogError(multiPkg.ErrTimeOut)
		return models.DailyActivity{}, multiPkg.ErrTimeOut
	}
}

func (c *MutexStorage) ReadToday(ctx context.Context, chatID int64) (map[string]models.DailyActivity, error) {
	select {
	case c.mutex <- true:
		defer func() {
			<-c.mutex
		}()
		return c.storage.ReadToday(ctx, chatID)
	case <-ctx.Done():
		// не дождались мьютекса
		journal.LogError(multiPkg.ErrTimeOut)
		return nil, multiPkg.ErrTimeOut
	}
}

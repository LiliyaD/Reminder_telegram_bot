package local

import (
	"context"
	"fmt"
	"time"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	cachePkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache"
	multiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/multithread"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	"github.com/pkg/errors"
)

type InterfaceStorage cachePkg.Interface

type cache struct {
	data  map[int64]map[string]models.DailyActivity //map[chatID}: map[activityName]: args
	chats map[int64]string                          // map[chatID}: userName
}

func newDataStorage() InterfaceStorage {
	return &cache{
		data:  map[int64]map[string]models.DailyActivity{},
		chats: map[int64]string{},
	}
}

func NewDataStorage() InterfaceStorage {
	return NewMutexStorage(newDataStorage())
}

func (c *cache) List(ctx context.Context, chatID int64, _ models.Pagination) (map[string]models.DailyActivity, error) {
	res := make(chan bool)
	go func() {
		_, ok := c.data[chatID]
		if ok {
			res <- true
		} else {
			res <- false
		}
	}()

	select {
	case a := <-res:
		// успех
		if a {
			return c.data[chatID], nil
		}
		journal.LogError(multiPkg.ErrEmptyStorage)
		return nil, multiPkg.ErrEmptyStorage
	case <-ctx.Done():
		// не дождались выполнения работы
		journal.LogError(multiPkg.ErrTimeOut)
		return nil, multiPkg.ErrTimeOut
	}
}

func (c *cache) ListStream(ctx context.Context, chatID int64, ch chan models.DailyActivityRec, chErr chan error) {
	res := make(chan bool)
	defer close(ch)
	defer ctx.Done()
	go func() {
		actM, ok := c.data[chatID]
		if ok {
			for actName, act := range actM {
				ch <- models.DailyActivityRec{
					Name:          actName,
					DailyActivity: act}
			}
			res <- true
		} else {
			res <- false
		}
	}()

	select {
	case a := <-res:
		// успех
		if a {
			chErr <- multiPkg.ErrOk
		}
		journal.LogError(multiPkg.ErrEmptyStorage)
		chErr <- multiPkg.ErrEmptyStorage
	case <-ctx.Done():
		// не дождались выполнения работы
		journal.LogError(multiPkg.ErrTimeOut)
		chErr <- multiPkg.ErrTimeOut
	}
}

func (c *cache) Add(ctx context.Context, actName string, act models.DailyActivity, chat models.Chat) (models.DailyActivity, error) {
	res := make(chan bool)
	go func() {
		if _, ok := c.initialRead(actName, chat.ChatID); ok {
			res <- false
		} else {
			actM, ok := c.data[chat.ChatID]
			if ok {
				_, ok := actM[actName]
				if ok {
					res <- false
				} else {
					actM[actName] = act
				}
			} else {
				actM := make(map[string]models.DailyActivity)
				actM[actName] = act
				c.data[chat.ChatID] = actM
			}

			c.chats[chat.ChatID] = chat.UserName
			res <- true
		}
	}()

	select {
	case a := <-res:
		// успех
		if a {
			return act, nil
		}
		journal.LogError(multiPkg.ErrExist)
		return models.DailyActivity{}, multiPkg.ErrExist
	case <-ctx.Done():
		// не дождались выполнения работы
		journal.LogError(multiPkg.ErrTimeOut)
		return models.DailyActivity{}, multiPkg.ErrTimeOut
	}
}

func (c *cache) Update(ctx context.Context, actName string, newAct models.DailyActivity, chatID int64) (models.DailyActivity, error) {
	var err error
	res := make(chan bool)
	var act models.DailyActivity
	go func() {
		var ok bool
		act, ok = c.initialRead(actName, chatID)
		if !ok {
			err = multiPkg.ErrNotExist
			res <- false
			return
		}

		prevBeginDate := act.BeginDate
		prevEndDate := act.EndDate
		t := time.Time{}
		if newAct.BeginDate != t {
			act.BeginDate = newAct.BeginDate
		}
		if newAct.EndDate != t {
			act.EndDate = newAct.EndDate
		}
		if act.EndDate.Before(act.BeginDate) {
			act.BeginDate = prevBeginDate
			act.EndDate = prevEndDate
			err = errors.New("End date is less than begin date")
			res <- false
			return
		}

		if newAct.TimesPerDay != 0 {
			act.TimesPerDay = newAct.TimesPerDay
		}

		if newAct.QuantityPerTime != 0 {
			act.QuantityPerTime = newAct.QuantityPerTime
		}

		actM, ok := c.data[chatID]
		if ok {
			_, ok := actM[actName]
			if !ok {
				res <- false
			} else {
				actM[actName] = act
			}
		} else {
			res <- false
		}

		res <- true
	}()

	select {
	case a := <-res:
		// успех
		if a {
			return act, nil
		}
		journal.LogError(err)
		return models.DailyActivity{}, err
	case <-ctx.Done():
		// не дождались выполнения работы
		journal.LogError(multiPkg.ErrTimeOut)
		return models.DailyActivity{}, multiPkg.ErrTimeOut
	}
}

func (c *cache) Delete(ctx context.Context, actName string, chatID int64) error {
	res := make(chan bool)
	go func() {
		if _, ok := c.initialRead(actName, chatID); !ok {
			res <- false
		} else {
			delete(c.data[chatID], actName)
			res <- true
		}
	}()

	select {
	case a := <-res:
		// успех
		if a {
			return nil
		}
		journal.LogError(multiPkg.ErrNotExist)
		return multiPkg.ErrNotExist
	case <-ctx.Done():
		// не дождались выполнения работы
		journal.LogError(multiPkg.ErrTimeOut)
		return multiPkg.ErrTimeOut
	}
}

func (c *cache) Read(ctx context.Context, actName string, chatID int64) (models.DailyActivity, error) {
	var act models.DailyActivity
	var ok bool
	res := make(chan bool)
	go func() {
		if act, ok = c.initialRead(actName, chatID); !ok {
			res <- false
		} else {
			res <- true
		}
	}()

	select {
	case a := <-res:
		// успех
		if a {
			return act, nil
		}
		journal.LogError(multiPkg.ErrNotExist)
		return models.DailyActivity{}, multiPkg.ErrNotExist
	case <-ctx.Done():
		// не дождались выполнения работы
		journal.LogError(multiPkg.ErrTimeOut)
		return models.DailyActivity{}, multiPkg.ErrTimeOut
	}
}

func (c *cache) ReadToday(ctx context.Context, chatID int64) (map[string]models.DailyActivity, error) {
	res := make(chan bool)
	actM := make(map[string]models.DailyActivity)
	go func() {
		act, ok := c.data[chatID]
		if ok {
			now := time.Now()
			d := fmt.Sprintf("%02d.%02d.%04d", now.Day(), now.Month(), now.Year())
			now, _ = time.Parse("02.01.2006", d)
			for key, val := range act {
				if (val.BeginDate.Before(now) && now.Before(val.EndDate)) || val.BeginDate == now || val.EndDate == now {
					actM[key] = val
				}
			}
			if len(actM) != 0 {
				res <- true
			} else {
				res <- false
			}
		} else {
			res <- false
		}
	}()

	select {
	case a := <-res:
		// успех
		if a {
			return actM, nil
		}
		journal.LogError(multiPkg.ErrNoActivitiesForToday)
		return nil, multiPkg.ErrNoActivitiesForToday
	case <-ctx.Done():
		// не дождались выполнения работы
		journal.LogError(multiPkg.ErrTimeOut)
		return nil, multiPkg.ErrTimeOut
	}
}

func (c *cache) initialRead(actName string, chatID int64) (models.DailyActivity, bool) {
	if act, ok := c.data[chatID]; ok {
		if act, ok := act[actName]; ok {
			return act, ok
		}
	}
	return models.DailyActivity{}, false
}

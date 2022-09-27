package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	cachePkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache"
	redis "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/database/cache"
	multiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/multithread"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/pkg/errors"
)

const (
	limit  = "ALL"
	offset = 0
	order  = "end_date DESC"
)

var Redis redis.InterfaceRedis

type InterfaceStorage interface {
	cachePkg.Interface
	CloseDataStorage()
}

type cache struct {
	pool  pgxpoolmock.PgxPool
	redis redis.InterfaceRedis
}

func NewDataStorage() InterfaceStorage {
	Redis = redis.Connect()
	return &cache{
		pool:  connect(),
		redis: Redis,
	}
}

func NewDataStoragePrep(pool pgxpoolmock.PgxPool) InterfaceStorage {
	return &cache{
		pool: pool,
	}
}

func GetCache() redis.InterfaceRedis {
	return Redis
}

func (c *cache) CloseDataStorage() {
	c.pool.Close()
}

func (c *cache) List(ctx context.Context, chatID int64, pagination models.Pagination) (map[string]models.DailyActivity, error) {
	var pag string
	if len(pagination.Order) == 0 {
		pag = "ORDER BY " + order
	} else {
		if orderF, ok := fieldsForOrder[pagination.Order]; ok {
			pag = "ORDER BY " + orderF
		} else {
			journal.LogError(multiPkg.ErrInvalidParameter)
			return nil, multiPkg.ErrInvalidParameter
		}
	}
	if pagination.Limit == 0 {
		pag += " LIMIT " + limit
	} else {
		pag += " LIMIT " + strconv.FormatUint(pagination.Limit, 10)
	}
	if pagination.Offset == 0 {
		pag += " OFFSET " + strconv.Itoa(offset)
	} else {
		pag += " OFFSET " + strconv.FormatUint(pagination.Limit, 10)
	}

	query := "SELECT act_name, begin_date, end_date, times_per_day, quantity_per_time FROM activities WHERE chat_id = $1 " + pag
	rows, err := c.pool.Query(ctx, query, chatID)
	if err != nil {
		journal.LogError(multiPkg.ErrInvalidParameter)
		return nil, multiPkg.ErrInvalidParameter
	}

	defer rows.Close()

	activities := make(map[string]models.DailyActivity)

	for rows.Next() {
		var act models.DailyActivity
		var actName string
		err := rows.Scan(&actName, &act.BeginDate, &act.EndDate, &act.TimesPerDay, &act.QuantityPerTime)
		if err != nil {
			journal.LogError(err)
			return nil, err
		}
		activities[actName] = act
	}

	if len(activities) == 0 {
		journal.LogError(multiPkg.ErrEmptyStorage)
		return nil, multiPkg.ErrEmptyStorage
	}

	return activities, nil
}

func (c *cache) ListStream(ctx context.Context, chatID int64, ch chan models.DailyActivityRec, chErr chan error) {
	query := "SELECT act_name, begin_date, end_date, times_per_day, quantity_per_time FROM activities WHERE chat_id = $1 "

	rows, err := c.pool.Query(ctx, query, chatID)
	if err != nil {
		journal.LogError(multiPkg.ErrInvalidParameter)
		chErr <- multiPkg.ErrInvalidParameter
		return
	}

	defer rows.Close()

	empty := true

	for rows.Next() {
		var act models.DailyActivity
		var actName string
		err := rows.Scan(&actName, &act.BeginDate, &act.EndDate, &act.TimesPerDay, &act.QuantityPerTime)
		if err != nil {
			journal.LogError(err)
			chErr <- err
		}

		empty = false

		ch <- models.DailyActivityRec{
			Name:          actName,
			DailyActivity: act}
	}

	if empty {
		journal.LogError(multiPkg.ErrEmptyStorage)
		chErr <- multiPkg.ErrEmptyStorage
	} else {
		chErr <- multiPkg.ErrOk
	}

}

func (c *cache) Add(ctx context.Context, actName string, act models.DailyActivity, chat models.Chat) (models.DailyActivity, error) {
	if c.redis != nil {
		fmt.Println("read from redis")
		var resp *models.DailyActivity
		if err := c.redis.Get(strconv.FormatInt(chat.ChatID, 10)+actName, &resp); err == nil {
			journal.LogError(multiPkg.ErrExist)
			return models.DailyActivity{}, multiPkg.ErrExist
		}
	}

	if _, ok := c.initialRead(ctx, actName, chat.ChatID); ok {
		journal.LogError(multiPkg.ErrExist)
		return models.DailyActivity{}, multiPkg.ErrExist
	}

	if err := c.addChat(ctx, chat); err != nil {
		journal.LogError(err)
		return models.DailyActivity{}, err
	}

	query := "INSERT INTO activities (chat_id, act_name, begin_date, end_date, times_per_day, quantity_per_time)" +
		"VALUES ($1,$2,$3,$4,$5,$6)"

	if _, err := c.pool.Exec(ctx, query,
		chat.ChatID, strings.ToLower(actName), act.BeginDate, act.EndDate, act.TimesPerDay, act.QuantityPerTime); err != nil {
		journal.LogError(err)
		return models.DailyActivity{}, err
	}

	return act, nil
}

func (c *cache) Update(ctx context.Context, actName string, newAct models.DailyActivity, chatID int64) (models.DailyActivity, error) {
	var act models.DailyActivity
	var ok bool
	act, ok = c.initialRead(ctx, actName, chatID)
	if !ok {
		journal.LogError(multiPkg.ErrNotExist)
		return models.DailyActivity{}, multiPkg.ErrNotExist
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
		err := errors.New("End date is less than begin date")
		journal.LogError(err)
		return models.DailyActivity{}, err
	}

	if newAct.TimesPerDay != 0 {
		act.TimesPerDay = newAct.TimesPerDay
	}

	if newAct.QuantityPerTime != 0 {
		act.QuantityPerTime = newAct.QuantityPerTime
	}

	query := "UPDATE activities SET begin_date = $1, end_date = $2, times_per_day = $3, quantity_per_time = $4" +
		"WHERE chat_id = $5 AND act_name = $6"

	if _, err := c.pool.Exec(ctx, query, act.BeginDate, act.EndDate, act.TimesPerDay, act.QuantityPerTime, chatID, strings.ToLower(actName)); err != nil {
		journal.LogError(err)
		return models.DailyActivity{}, err
	}

	if c.redis != nil {
		c.redis.Del(strconv.FormatInt(chatID, 10) + actName)
	}

	return act, nil
}

func (c *cache) Delete(ctx context.Context, actName string, chatID int64) error {
	if _, ok := c.initialRead(ctx, actName, chatID); !ok {
		journal.LogError(multiPkg.ErrNotExist)
		return multiPkg.ErrNotExist
	}

	query := "DELETE FROM activities WHERE chat_id = $1 AND act_name = $2"

	if _, err := c.pool.Exec(ctx, query, chatID, strings.ToLower(actName)); err != nil {
		journal.LogError(err)
		return err
	}

	if c.redis != nil {
		c.redis.Del(strconv.FormatInt(chatID, 10) + actName)
	}

	return nil
}

func (c *cache) Read(ctx context.Context, actName string, chatID int64) (models.DailyActivity, error) {
	if c.redis != nil {
		fmt.Println("read from redis")
		var resp *models.DailyActivity
		if err := c.redis.Get(strconv.FormatInt(chatID, 10)+actName, &resp); err == nil {
			return *resp, nil
		}
	}

	if act, ok := c.initialRead(ctx, actName, chatID); !ok {
		journal.LogError(multiPkg.ErrNotExist)
		return models.DailyActivity{}, multiPkg.ErrNotExist
	} else {
		if c.redis != nil {
			fmt.Println("write to redis")
			c.redis.Set(strconv.FormatInt(chatID, 10)+actName, act, 0)
		}
		return act, nil
	}
}

func (c *cache) ReadToday(ctx context.Context, chatID int64) (map[string]models.DailyActivity, error) {
	query := "SELECT act_name, begin_date, end_date, times_per_day, quantity_per_time FROM activities WHERE chat_id = $1 AND begin_date <= current_date AND end_date >= current_date"

	rows, err := c.pool.Query(ctx, query, chatID)
	if err != nil {
		journal.LogError(multiPkg.ErrInvalidParameter)
		return nil, multiPkg.ErrInvalidParameter
	}

	defer rows.Close()

	activities := make(map[string]models.DailyActivity)

	for rows.Next() {
		var act models.DailyActivity
		var actName string
		err := rows.Scan(&actName, &act.BeginDate, &act.EndDate, &act.TimesPerDay, &act.QuantityPerTime)
		if err != nil {
			journal.LogError(err)
			return nil, err
		}
		activities[actName] = act
	}

	if len(activities) == 0 {
		journal.LogError(multiPkg.ErrNoActivitiesForToday)
		return nil, multiPkg.ErrNoActivitiesForToday
	}

	return activities, nil
}

func (c *cache) initialRead(ctx context.Context, actName string, chatID int64) (models.DailyActivity, bool) {
	var act models.DailyActivity
	ok := true
	query := "SELECT begin_date, end_date, times_per_day, quantity_per_time FROM activities WHERE chat_id = $1 AND act_name = $2"
	row := c.pool.QueryRow(ctx, query, chatID, strings.ToLower(actName))
	if err := row.Scan(&act.BeginDate, &act.EndDate, &act.TimesPerDay, &act.QuantityPerTime); err != nil {
		journal.LogInfo(err)
		ok = false
	}
	return act, ok
}

func (c *cache) addChat(ctx context.Context, chat models.Chat) error {
	query := "SELECT chat_id FROM chats WHERE chat_id = $1"
	row := c.pool.QueryRow(ctx, query, chat.ChatID)
	var chatID int64
	err := row.Scan(&chatID)
	if err != nil {
		query = "INSERT INTO chats (chat_id, user_name) VALUES ($1,$2)"
		_, err = c.pool.Exec(ctx, query, chat.ChatID, chat.UserName)
		if err != nil {
			journal.LogError(err)
			return err
		}
	}
	return nil
}

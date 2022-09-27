package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	"github.com/edbmanniwood/pgxpoolmock"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ = (func() interface{} {
	_testing = true
	return nil
}())

func TestNewActivity(t *testing.T) {
	t.Run("Success db creation new chat, activity", func(t *testing.T) {
		//arrange
		f := connectMock(t)
		defer f.closeMockConnection()
		ctx := context.Background()

		beginDate, err := time.Parse("02.01.2006", "10.07.2022")
		if err != nil {
			t.Fatal(err)
		}
		endDate, err := time.Parse("02.01.2006", "15.09.2023")
		if err != nil {
			t.Fatal(err)
		}

		act := models.DailyActivity{
			BeginDate:       beginDate,
			EndDate:         endDate,
			TimesPerDay:     2,
			QuantityPerTime: 30}

		chat := models.Chat{
			ChatID:   777,
			UserName: "agent007"}

		actName := "push-up"

		// tables definition
		columnsAct := []string{"chat_id", "act_name", "begin_date", "end_date", "times_per_day", "quantity_per_time"}
		columnsChat := []string{"chat_id", "user_name"}

		query1 := "SELECT begin_date, end_date, times_per_day, quantity_per_time FROM activities WHERE chat_id = $1 AND act_name = $2"
		returnedRow1 := pgxpoolmock.NewRows(columnsAct).ToPgxRows()
		f.mockPool.EXPECT().QueryRow(ctx, query1, chat.ChatID, actName).Return(returnedRow1)

		query2 := "SELECT chat_id FROM chats WHERE chat_id = $1"
		returnedRow2 := pgxpoolmock.NewRows(columnsChat).ToPgxRows()
		f.mockPool.EXPECT().QueryRow(ctx, query2, chat.ChatID).Return(returnedRow2)

		query3 := "INSERT INTO chats (chat_id, user_name) VALUES ($1,$2)"
		commandTag1 := pgconn.CommandTag("INSERT 0 1")
		f.mockPool.EXPECT().Exec(ctx, query3, chat.ChatID, chat.UserName).Return(commandTag1, nil)

		query4 := "INSERT INTO activities (chat_id, act_name, begin_date, end_date, times_per_day, quantity_per_time)" +
			"VALUES ($1,$2,$3,$4,$5,$6)"
		commandTag2 := pgconn.CommandTag("INSERT 0 1 2 3 4 5")
		f.mockPool.EXPECT().Exec(ctx, query4, chat.ChatID, actName, act.BeginDate, act.EndDate, act.TimesPerDay, act.QuantityPerTime).Return(commandTag2, nil)

		// act
		result, err := f.storage.Add(ctx, actName, act, chat)

		//assert
		require.NoError(t, err)
		assert.Equal(t, act, result)
	})
}

func TestReadActivity(t *testing.T) {
	t.Run("Success db read activity", func(t *testing.T) {
		//arrange
		f := connectMock(t)
		defer f.closeMockConnection()
		ctx := context.Background()

		chatID := int64(777)
		actName := "push-up"

		beginDate, err := time.Parse("02.01.2006", "10.07.2022")
		if err != nil {
			t.Fatal(err)
		}
		endDate, err := time.Parse("02.01.2006", "15.09.2023")
		if err != nil {
			t.Fatal(err)
		}

		act := models.DailyActivity{
			BeginDate:       beginDate,
			EndDate:         endDate,
			TimesPerDay:     2,
			QuantityPerTime: 30}

		query1 := "SELECT begin_date, end_date, times_per_day, quantity_per_time FROM activities WHERE chat_id = $1 AND act_name = $2"

		columns := []string{"begin_date", "end_date", "times_per_day", "quantity_per_time"}

		returnedRow := pgxpoolmock.NewRow(columns, act.BeginDate, act.EndDate, act.TimesPerDay, act.QuantityPerTime)
		f.mockPool.EXPECT().QueryRow(ctx, query1, chatID, actName).Return(returnedRow)

		// act
		result, err := f.storage.Read(ctx, actName, chatID)

		//assert
		require.NoError(t, err)
		assert.Equal(t, act, result)
	})
}

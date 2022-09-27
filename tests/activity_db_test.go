//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"

	multiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/multithread"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBCreateActivity(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Activity was created", func(t *testing.T) {
			//arrange in init_test.go

			//act
			resp, err := Storage.Add(context.Background(), actName, act, chat)

			//assert
			require.NoError(t, err)
			assert.Equal(t, resp, act)
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("AlreadyExists: activity exists", func(t *testing.T) {
			//arrange in init_test.go

			//act
			resp, err := Storage.Add(context.Background(), actName, act, chat)

			//assert
			assert.EqualError(t, err, multiPkg.ErrExist.Error())
			assert.Empty(t, resp)
		})
	})
}

func TestDBReadActivity(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Activity was found", func(t *testing.T) {
			//arrange in init_test.go
			Storage.Add(context.Background(), actName, act, chat)

			//act
			resp, err := Storage.Read(context.Background(), actName, chat.ChatID)

			//assert
			require.NoError(t, err)
			assert.Equal(t, resp, act)
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("ErrNotExist: activity doesn't exist", func(t *testing.T) {
			//arrange in init_test.go

			//act
			resp, err := Storage.Read(context.Background(), "pushup", chat.ChatID)

			//assert
			require.EqualError(t, err, multiPkg.ErrNotExist.Error())
			assert.Empty(t, resp)
		})
	})
}

func TestDBListActivity(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Activity was created", func(t *testing.T) {
			//arrange in init_test.go
			Storage.Add(context.Background(), actName, act, chat)

			//act
			resp, err := Storage.List(context.Background(), chat.ChatID, pagination)

			//assert
			require.NoError(t, err)
			assert.Equal(t, resp, activitiesMap)
		})
	})
}

func TestDBUpdateActivity(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Activity was created", func(t *testing.T) {
			//arrange in init_test.go
			Storage.Add(context.Background(), actName, act, chat)

			//act
			resp, err := Storage.Update(context.Background(), actName, models.DailyActivity{TimesPerDay: 3}, chat.ChatID)

			//assert
			require.NoError(t, err)
			assert.Equal(t, resp, models.DailyActivity{
				BeginDate:       beginDate,
				EndDate:         endDate,
				TimesPerDay:     3,
				QuantityPerTime: 30})
		})
	})
}

func TestDBDeleteActivity(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Activity was deleted", func(t *testing.T) {
			//arrange in init_test.go
			Storage.Add(context.Background(), actName, act, chat)

			//act
			err := Storage.Delete(context.Background(), actName, chat.ChatID)

			//assert
			require.NoError(t, err)
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("ErrNotExist: activity doesn't exist", func(t *testing.T) {
			//arrange in init_test.go

			//act
			err := Storage.Delete(context.Background(), actName, chat.ChatID)

			//assert
			require.EqualError(t, err, multiPkg.ErrNotExist.Error())
		})
	})
}
